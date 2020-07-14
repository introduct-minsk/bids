--
-- Data initialization
--
begin;

insert into bid_source
values ('game'),
       ('server'),
       ('payment');

insert into player(phone_number, balance)
values ('1', 0);
commit;