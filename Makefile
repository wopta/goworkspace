generate_partnership_key:
	openssl rand -base64 32 # 32 is the number of bytes

.PHONE: generate_partnership_key
.SILENT: generate_partnership_key
