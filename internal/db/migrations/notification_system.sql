-- `notification_system.sql` --
-- CLEANING_STRUCTURED --
drop view if exists push_relay_events_context;

drop table if exists public.new_context;

drop table if exists public.push_relay_events;

drop trigger if exists push_relay_trigger_status_change on tb_atend;

drop trigger if exists trigger_status_change on tb_atend;

drop trigger if exists push_relay_trigger_status_change on tb_atend_prof;

drop trigger if exists trigger_status_change on tb_atend_prof;

drop function if exists push_relay_notify_status_change();
	
-- function
create or replace function push_relay_notify_status_change() returns trigger as $$
declare
    payload JSONB;

begin
    if NEW.tp_atend_prof > 0 THEN
		select row_to_json(sub) payload
		into payload
		from (select
				t1.co_seq_atend,
				tp2.nu_cns as prof_cns,
				tc2.co_cbo_2002 as prof_cbo,
				tus.nu_cnes as cnes ,
				ttap.no_tipo_atend_prof as local_chamada,
				UPPER(tc.no_cidadao) as cidadao
			from
				tb_atend t1 join tb_atend_prof tap on t1.co_seq_atend = tap.co_atend
			join tb_lotacao tl on tap.co_lotacao = tl.co_ator_papel
			join tb_cbo tc2 on tl.co_cbo = tc2.co_cbo
			join tb_prontuario tp on t1.co_prontuario = tp.co_seq_prontuario
			join tb_cidadao tc on tp.co_cidadao = tc.co_seq_cidadao
			join tb_unidade_saude tus on t1.co_unidade_saude = tus.co_seq_unidade_saude
			join tb_status_atend tsa on t1.st_atend = tsa.co_status_atend
			join tb_prof tp2 on tl.co_prof = tp2.co_seq_prof
			join tb_tipo_atend_prof ttap on tap.tp_atend_prof = ttap.co_tipo_atend_prof
			where t1.co_seq_atend = NEW.co_atend
			AND to_char(t1.dt_criacao_registro,'dd/mm/yyyy') = to_char(now(),'dd/mm/yyyy')
			limit 1
		) AS sub;
		perform pg_notify('call_record', payload::text);
		return NEW;
	else
		return OLD;
	end if;
end;
$$ language plpgsql;

-- trigger
CREATE TRIGGER push_relay_trigger_status_change
    AFTER update 
    ON public.tb_atend_prof
    FOR EACH ROW
    EXECUTE PROCEDURE public.push_relay_notify_status_change();