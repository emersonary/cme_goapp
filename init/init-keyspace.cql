CREATE KEYSPACE IF NOT EXISTS cme
WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 3};