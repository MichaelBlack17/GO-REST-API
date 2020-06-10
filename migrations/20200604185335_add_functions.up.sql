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