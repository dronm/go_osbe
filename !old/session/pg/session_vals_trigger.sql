    CREATE TRIGGER session_vals_trigger_after
    AFTER DELETE
    ON public.session_vals
    FOR EACH ROW
    EXECUTE PROCEDURE public.session_vals_process();
