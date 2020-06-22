--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2
-- Dumped by pg_dump version 12.2

-- Started on 2020-06-22 12:28:38

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 2 (class 3079 OID 65614)
-- Name: pldbgapi; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pldbgapi WITH SCHEMA public;


--
-- TOC entry 2915 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pldbgapi; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pldbgapi IS 'server-side support for debugging PL/pgSQL functions';


--
-- TOC entry 218 (class 1255 OID 65583)
-- Name: addrequest(bigint, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.addrequest(userid bigint, text character varying) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
DECLARE myid bigint;
BEGIN
INSERT INTO public.requests(
	 user_id, message, create_date)
	VALUES (userid, text, clock_timestamp()) RETURNING id INTO myid;
Return myid; 
END
$$;


ALTER FUNCTION public.addrequest(userid bigint, text character varying) OWNER TO postgres;

--
-- TOC entry 233 (class 1255 OID 65587)
-- Name: addrequestinqueue(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.addrequestinqueue() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE MngId bigint;
DECLARE queuelen int = 2;
DECLARE validtime int = 15;
BEGIN
	 SELECT glb('queue_length')  INTO queuelen;
	 SELECT glb('valid_time')  INTO validtime;
	
	SELECT id INTO MngId 
	FROM 
		public.managers
	EXCEPT
	SELECT mid.manager_id 
	FROM (
	SELECT manager_id, COUNT(manager_id) AS cnt FROM public.requestqueue
	GROUP BY manager_id) mid WHERE mid.cnt > (queuelen - 1) 
	LIMIT 1 ;
			 	
	--WHERE COALESCE(array_length("Queue", 1), 0) < 2	
	--ORDER BY COALESCE(array_length("Queue", 1), 0) ASC
	
	
	IF MngId > 0  
	THEN
	INSERT INTO public.RequestQueue(request_id,status, manager_id, valid_time) VALUES
	(New.id, 1, MngId, (CURRENT_TIMESTAMP + (validtime * interval '1 minute')));
	
	ELSE
		INSERT INTO public.RequestQueue(request_id, status) VALUES
	(New.Id,0);
	END IF;

	
	RETURN new;
END;
$$;


ALTER FUNCTION public.addrequestinqueue() OWNER TO postgres;

--
-- TOC entry 217 (class 1255 OID 65611)
-- Name: cancelprocessingrequest(bigint, bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.cancelprocessingrequest(mngid bigint, reqid bigint) RETURNS json
    LANGUAGE plpgsql
    AS $$
	DECLARE myrow public.requestqueue%rowtype;
BEGIN
	UPDATE public.requestqueue SET status = 0, valid_time = NULL, manager_id = NULL WHERE request_id = reqid AND manager_id = mngid RETURNING * INTO myrow;
	return row_to_json(myrow);
END
$$;


ALTER FUNCTION public.cancelprocessingrequest(mngid bigint, reqid bigint) OWNER TO postgres;

--
-- TOC entry 234 (class 1255 OID 65599)
-- Name: cancelrequest(bigint, bigint); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.cancelrequest(usrid bigint, reqid bigint) RETURNS json
    LANGUAGE plpgsql
    AS $$
	DECLARE myrow public.requestqueue%rowtype;
BEGIN
	
	DELETE FROM public.requestqueue WHERE id IN
	(select rq.id FROM public.requestqueue rq
	JOIN public.requests r ON rq.request_id = r.id
	WHERE (rq.request_id = reqid) AND (r.user_id = usrid)) returning * INTO myrow;
	
	return row_to_json(myrow);
END
$$;


ALTER FUNCTION public.cancelrequest(usrid bigint, reqid bigint) OWNER TO postgres;

--
-- TOC entry 219 (class 1255 OID 65612)
-- Name: glb(text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.glb(code text) RETURNS integer
    LANGUAGE sql
    AS $$
    select current_setting('glb.' || code)::integer;
$$;


ALTER FUNCTION public.glb(code text) OWNER TO postgres;

--
-- TOC entry 255 (class 1255 OID 65613)
-- Name: querymanagement(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.querymanagement() RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE doneids bigint[];
DECLARE lateids bigint[];
DECLARE ids bigint[];
DECLARE mngrs bigint[];
DECLARE managerid bigint;
DECLARE m bigint;
DECLARE indx int;
DECLARE QueueLen int = 2;
DECLARE validtime int = 15;
BEGIN
	set glb.queue_length to 3;
	set glb.valid_time to 1;
	select glb('queue_length')  INTO QueueLen;
	select glb('valid_time')  INTO validtime;
	
	--массив отработанных id
	doneids := ARRAY(
	SELECT request_id FROM public.requestqueue
	WHERE status = 2
	);
	
	--удаляем отработанные заявки из очереди
	DELETE FROM public.requestqueue
	WHERE array_position(doneids, request_id) IS NOT NULL;
	
	--получаем массив Id которые висят на менеджере больше валидного времени или ожидают
	ids := ARRAY(SELECT request_id FROM public.requestqueue
	WHERE (status = 1 and valid_time < CURRENT_TIMESTAMP) OR (status = 0) ORDER BY status DESC);
	
	--получаем массив доступных id менеджеров
	mngrs := ARRAY(
		(SELECT t.id
			FROM (select id, generate_series(1, QueueLen) FROM public.managers) t)
		EXCEPT ALL
			(SELECT manager_id
			FROM public.requestqueue
			WHERE status = 1 AND valid_time > CURRENT_TIMESTAMP));
	
	indx = 0;	
  FOREACH m IN ARRAY ids
  LOOP
  	managerid = 0;
	IF indx < array_length(mngrs,1) THEN
		indx = indx + 1;
		managerid = mngrs[indx];
	END IF;
	IF managerid > 0  
	THEN
		UPDATE public.requestqueue
		SET manager_id = managerid,
		status = 1,
		valid_time = (CURRENT_TIMESTAMP + (validtime * interval '1 minute')) 
		WHERE request_id = m;
	ELSE
		UPDATE public.requestqueue 
		SET manager_id = NULL,
		status = 0,
		valid_time = NULL
		WHERE request_id = m;

	END IF;
  END LOOP;
	
return 0;
END
$$;


ALTER FUNCTION public.querymanagement() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 207 (class 1259 OID 57491)
-- Name: managers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.managers (
    id bigint NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.managers OWNER TO postgres;

--
-- TOC entry 206 (class 1259 OID 57489)
-- Name: managers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.managers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.managers_id_seq OWNER TO postgres;

--
-- TOC entry 2916 (class 0 OID 0)
-- Dependencies: 206
-- Name: managers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.managers_id_seq OWNED BY public.managers.id;


--
-- TOC entry 211 (class 1259 OID 57521)
-- Name: requestqueue; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.requestqueue (
    id bigint NOT NULL,
    request_id bigint NOT NULL,
    manager_id bigint,
    status integer NOT NULL,
    valid_time timestamp with time zone
);


ALTER TABLE public.requestqueue OWNER TO postgres;

--
-- TOC entry 210 (class 1259 OID 57519)
-- Name: requestqueue_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.requestqueue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.requestqueue_id_seq OWNER TO postgres;

--
-- TOC entry 2917 (class 0 OID 0)
-- Dependencies: 210
-- Name: requestqueue_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.requestqueue_id_seq OWNED BY public.requestqueue.id;


--
-- TOC entry 209 (class 1259 OID 57506)
-- Name: requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.requests (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    message character varying NOT NULL,
    create_date timestamp with time zone NOT NULL
);


ALTER TABLE public.requests OWNER TO postgres;

--
-- TOC entry 208 (class 1259 OID 57504)
-- Name: requests_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.requests_id_seq OWNER TO postgres;

--
-- TOC entry 2918 (class 0 OID 0)
-- Dependencies: 208
-- Name: requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.requests_id_seq OWNED BY public.requests.id;


--
-- TOC entry 203 (class 1259 OID 49195)
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- TOC entry 205 (class 1259 OID 57476)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 57474)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 2919 (class 0 OID 0)
-- Dependencies: 204
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 2761 (class 2604 OID 57494)
-- Name: managers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.managers ALTER COLUMN id SET DEFAULT nextval('public.managers_id_seq'::regclass);


--
-- TOC entry 2763 (class 2604 OID 57524)
-- Name: requestqueue id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.requestqueue ALTER COLUMN id SET DEFAULT nextval('public.requestqueue_id_seq'::regclass);


--
-- TOC entry 2762 (class 2604 OID 57509)
-- Name: requests id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.requests ALTER COLUMN id SET DEFAULT nextval('public.requests_id_seq'::regclass);


--
-- TOC entry 2760 (class 2604 OID 57479)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 2905 (class 0 OID 57491)
-- Dependencies: 207
-- Data for Name: managers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.managers (id, name) FROM stdin;
1	manager 1
2	manager 2
\.


--
-- TOC entry 2909 (class 0 OID 57521)
-- Dependencies: 211
-- Data for Name: requestqueue; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.requestqueue (id, request_id, manager_id, status, valid_time) FROM stdin;
\.


--
-- TOC entry 2907 (class 0 OID 57506)
-- Dependencies: 209
-- Data for Name: requests; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.requests (id, user_id, message, create_date) FROM stdin;
\.


--
-- TOC entry 2901 (class 0 OID 49195)
-- Dependencies: 203
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20200604185335	f
\.


--
-- TOC entry 2903 (class 0 OID 57476)
-- Dependencies: 205
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name) FROM stdin;
1	Test
2	Test
3	Test
4	Test
5	Vasia
\.


--
-- TOC entry 2920 (class 0 OID 0)
-- Dependencies: 206
-- Name: managers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.managers_id_seq', 2, true);


--
-- TOC entry 2921 (class 0 OID 0)
-- Dependencies: 210
-- Name: requestqueue_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.requestqueue_id_seq', 105, true);


--
-- TOC entry 2922 (class 0 OID 0)
-- Dependencies: 208
-- Name: requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.requests_id_seq', 124, true);


--
-- TOC entry 2923 (class 0 OID 0)
-- Dependencies: 204
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 5, true);


--
-- TOC entry 2769 (class 2606 OID 57499)
-- Name: managers managers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.managers
    ADD CONSTRAINT managers_pkey PRIMARY KEY (id);


--
-- TOC entry 2773 (class 2606 OID 57526)
-- Name: requestqueue requestqueue_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.requestqueue
    ADD CONSTRAINT requestqueue_pkey PRIMARY KEY (id);


--
-- TOC entry 2771 (class 2606 OID 57514)
-- Name: requests requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_pkey PRIMARY KEY (id);


--
-- TOC entry 2765 (class 2606 OID 49199)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 2767 (class 2606 OID 57484)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 2774 (class 2620 OID 65588)
-- Name: requests queue_insert; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER queue_insert AFTER INSERT ON public.requests FOR EACH ROW EXECUTE FUNCTION public.addrequestinqueue();


-- Completed on 2020-06-22 12:28:38

--
-- PostgreSQL database dump complete
--

