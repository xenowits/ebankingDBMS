create table transactions(
  id INT PRIMARY KEY AUTO_INCREMENT ,
  t_datetime VARCHAR(255) NOT NULL,
  user_credited VARCHAR(255) NOT NULL,
  user_debited VARCHAR(255) NOT NULL,
  creditedUser_finalBalance INT DEFAULT 0,
  debitedUser_finalBalance INT DEFAULT 0,
  t_amount INT NOT NULL
);

select id, user_credited, user_debited, creditedUser_finalBalance, debitedUser_finalBalance, t_amount from transactions;

create table customers(
  id INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(255),
  availBalance INT DEFAULT 0
);

select * from customers;

insert into customers (username, password, role, availBalance) VALUES (
"abhishek",
"$2a$05$ZZVw36KZsg8zLNpq2jeX9eTxOXDh.wGWybDMtKDCP7aqDdbBI4ukW",
"admin",
0);
