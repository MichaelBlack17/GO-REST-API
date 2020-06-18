CREATE OR REPLACE FUNCTION public.addrequest(
    userid bigint,
    text character varying)
    RETURNS bigint
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE

AS $BODY$
DECLARE myid bigint;
BEGIN
    INSERT INTO public.requests(
        userid, message, createdate)
    VALUES (userid, text, clock_timestamp()) RETURNING id INTO myid;
    Return myid;
END
$BODY$;

ALTER FUNCTION public.addrequest(character varying)
    OWNER TO postgres;

-------------------------------------
CREATE FUNCTION public.addrequestinqueue()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
DECLARE ManagerId bigint;
BEGIN
    SELECT id INTO ManagerId
    FROM
        public.managers
        EXCEPT
    SELECT mid.managerid
    FROM (
             SELECT managerid, COUNT(managerid) AS cnt FROM public.requestqueue
             GROUP BY managerid) mid WHERE mid.cnt < 2
    LIMIT 1 ;

    --WHERE COALESCE(array_length("Queue", 1), 0) < 2
    --ORDER BY COALESCE(array_length("Queue", 1), 0) ASC


    IF ManagerId > 0
    THEN
        INSERT INTO public.RequestQueue(requestid,status, managerid, validtime) VALUES
        (New.id, 1, ManagerId, (CURRENT_TIMESTAMP + (15 * interval '1 minute')));

    ELSE
        INSERT INTO public.RequestQueue(requestid, status) VALUES
        (New.Id,0);
    END IF;


    RETURN new;
END;
$BODY$;

ALTER FUNCTION public.addrequestinqueue()
    OWNER TO postgres;



CREATE TRIGGER queue_insert
    AFTER INSERT
    ON public.requests
    FOR ROW
EXECUTE PROCEDURE addrequestinqueue();
----------------------------------

CREATE OR REPLACE FUNCTION public.cancelrequest(
    usrid bigint,
    reqid bigint)
    RETURNS json
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE

AS $BODY$
DECLARE myrow public.requestqueue%rowtype;
BEGIN

    DELETE FROM public.requestqueue WHERE id IN
                                          (select rq.id FROM public.requestqueue rq
                                                                 JOIN public.requests r ON rq.request_id = r.id
                                           WHERE (rq.request_id = reqid) AND (r.user_id = usrid)) returning * INTO myrow;

    return row_to_json(myrow);
END
$BODY$;

ALTER FUNCTION public.cancelrequest(bigint, bigint)
    OWNER TO postgres;
--------------------------------------------------------------
CREATE OR REPLACE FUNCTION public.cancelprocessingrequest(
    mngid bigint,
    reqid bigint)
    RETURNS json
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE

AS $BODY$
DECLARE myrow public.requestqueue%rowtype;
BEGIN
    UPDATE public.requestqueue SET status = 0, valid_time = NULL, manager_id = NULL WHERE request_id = reqid AND manager_id = mngid RETURNING * INTO myrow;
    return row_to_json(myrow);
END
$BODY$;

ALTER FUNCTION public.cancelprocessingrequest(bigint, bigint)
    OWNER TO postgres;
---------------------------------------------------------------

create or replace function public.glb(code text)
    returns integer language sql as $$
select current_setting('glb.' || code)::integer;
$$;

ALTER FUNCTION public.glb(text)
    OWNER TO postgres;
------------------------------------------------------------------
CREATE OR REPLACE FUNCTION public.querymanagement(
)
    RETURNS integer
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE

AS $BODY$
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
$BODY$;

ALTER FUNCTION public.querymanagement()
    OWNER TO postgres;
