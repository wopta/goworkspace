generate_partnership_key:
	openssl rand -base64 32 # 32 is the number of bytes

test_trigger:
	act --container-architecture linux/amd64 --secret-file .secrets --rm

.SILENT: generate_partnership_key test_trigger
