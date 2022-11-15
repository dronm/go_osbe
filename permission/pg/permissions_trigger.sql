    CREATE TRIGGER permissions_trigger_after
    AFTER DELETE OR UPDATE OR INSERT
    ON public.permissions
    FOR EACH ROW
    EXECUTE PROCEDURE public.permissions_process();
