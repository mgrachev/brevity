CREATE SEQUENCE global_id_sequence;

CREATE OR REPLACE FUNCTION generate_primary_key(OUT RESULT BIGINT) AS
$$
DECLARE
    our_epoch BIGINT := 1314220021721;
    seq_id BIGINT;
    now_millis BIGINT;
    url_shortener_id INT := 1;
BEGIN
    SELECT nextval('global_id_sequence') % 1024 INTO seq_id;

    SELECT floor(EXTRACT(epoch FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23;
    result := RESULT | (url_shortener_id << 10);
    result := RESULT | (seq_id);
END;
$$
LANGUAGE plpgsql;
