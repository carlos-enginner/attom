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
SELECT
	t1.co_seq_atend,
	tts.co_seq_tipo_servico co_tipo_servico,
	tts.no_tipo_servico no_servico,
    t1.co_seq_atend AS senha,
    t2.no_status_atend AS status,
    t2.co_status_atend AS estagio,
    UPPER(t4.no_cidadao) AS cidadao,
    t4.dt_nascimento,
    t4.nu_cpf AS cpf,
    t4.nu_cns AS cns,
    CASE
        WHEN t4.no_sexo = 'FEMININO' THEN 'F'
        WHEN t4.no_sexo = 'MASCULINO' THEN 'M'
    END AS sexo,
    t1.dt_criacao_registro AS tempo,
    UPPER(t6.no_profissional) AS profissional,
    t6.nu_cns AS prof_cns,
    t6.nu_cpf AS prof_cpf,
    t8.co_cbo_2002 AS prof_cbo_nu,
    t8.no_cbo AS prof_cbo,
    t10.nu_ine AS prof_ine,
    t10.no_equipe AS prof_equipe,
    t9.nu_cnes AS cnes,
    t9.no_unidade_saude AS unidade,
    t1.dt_ultima_alteracao_status AS tempo_atendido,
    t1.dt_ultima_alteracao_status AS tempo_cancelado,
    t11.no_classificacao_risco,
    t11.co_classificacao_risco
FROM
    tb_atend t1
LEFT JOIN tb_status_atend t2 ON t2.co_status_atend = t1.st_atend
LEFT JOIN tb_prontuario t3 ON t3.co_seq_prontuario = t1.co_prontuario
LEFT JOIN tb_cidadao t4 ON t4.co_seq_cidadao = t3.co_cidadao
LEFT JOIN tb_ator_papel t5 ON t1.co_responsavel = t5.co_seq_ator_papel
LEFT JOIN tb_prof t6 ON t6.co_seq_prof = t5.co_prof
LEFT JOIN tb_lotacao t7 ON t7.co_prof = t6.co_seq_prof
LEFT JOIN tb_cbo t8 ON t8.co_cbo = t7.co_cbo
LEFT JOIN tb_unidade_saude t9 ON t1.co_unidade_saude = t9.co_seq_unidade_saude
LEFT JOIN tb_equipe t10 ON t1.co_equipe = t10.co_seq_equipe
LEFT JOIN tb_classificacao_risco t11 ON t1.co_classificacao_risco = t11.co_classificacao_risco
LEFT join rl_atend_tipo_servico rats ON	t1.co_seq_atend = rats.co_atend
join tb_tipo_servico tts ON rats.tp_servico = tts.co_seq_tipo_servico;


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