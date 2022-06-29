
list:
	@ cat Makefile

gin:
	gin -i -d cmd/server

psql:
	@ psql 'postgres://pguser:pgpassword@localhost/postgres'
