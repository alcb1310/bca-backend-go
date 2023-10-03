create view user_without_password as
select u.id as id, u.email as email, u.name as name, u.company_id as company_id, r.name as role
from "user" u inner join role r on u.role_id = r.id;

create view total_users as
select company_id company_id, count(*) total
from "user"
group by company_id;
