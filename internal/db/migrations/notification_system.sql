-- `notification_system.sql` --
-- CLEANING_STRUCTURED --
drop view if exists push_relay_events_context;

drop table if exists public.push_relay_events;

drop trigger if exists push_relay_trigger_status_change on tb_atend;

drop trigger if exists trigger_status_change on tb_atend;

drop function if exists push_relay_notify_status_change();

-- CREATING_STRUCTURED --
-- view
CREATE VIEW push_relay_events_context AS
select
	t1.co_seq_atend,
	tp2.nu_cns as prof_cns,
	tc2.co_cbo_2002 as prof_cbo,
	tus.nu_cnes as cnes ,
	trim(replace (tsa.no_status_atend, 'EM', '')) as local_chamada,
	UPPER(tc.no_cidadao) as cidadao
from
	tb_atend t1
join tb_atend_prof tap on
	t1.co_seq_atend = tap.co_atend
join tb_lotacao tl on
	tap.co_lotacao = tl.co_ator_papel
join tb_cbo tc2 on
	tl.co_cbo = tc2.co_cbo
join tb_prontuario tp on
	t1.co_prontuario = tp.co_seq_prontuario
join tb_cidadao tc on
	tp.co_cidadao = tc.co_seq_cidadao
join tb_unidade_saude tus on
	t1.co_unidade_saude = tus.co_seq_unidade_saude
join tb_status_atend tsa on
	t1.st_atend = tsa.co_status_atend
join tb_prof tp2 on
	tl.co_prof = tp2.co_seq_prof
where
	to_char(t1.dt_criacao_registro ,
	'dd/mm/yyyy') = to_char(now() ,
	'dd/mm/yyyy');


-- table
create table public.push_relay_events (
    id SERIAL primary key,
    status INT,
    context JSONB,
    updated_at TIMESTAMP default CURRENT_TIMESTAMP
);
	
-- function
create or replace function push_relay_notify_status_change() returns trigger as $$
declare
    new_context JSONB;

begin
    if NEW.st_atend in (2, 3) then
-- Create a JSON object to store in the context column
        select
	row_to_json(u)
into
	new_context
from
	push_relay_events_context u
where
	u.co_seq_atend = NEW.co_seq_atend LIMIT 1;
-- Insert a new row into the events table
       insert
	into
	push_relay_events (status,
	context)
values (NEW.st_atend, new_context);
-- Notify when status changes
-- PERFORM pg_notify('status_change', 'Status changed to ' || NEW.st_atend || ' for ID: ' || NEW.id);
perform pg_notify('status_change', new_context::text);
end if;

return new;
end;

$$ language plpgsql;

-- trigger
CREATE TRIGGER push_relay_trigger_status_change
    AFTER update 
    ON public.tb_atend
    FOR EACH ROW
    EXECUTE PROCEDURE public.push_relay_notify_status_change(); 