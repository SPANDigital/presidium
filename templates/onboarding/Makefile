default:
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info 

hugo:
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info 

tidy:
	go mod tidy
	hugo mod tidy

refresh:
	hugo mod clean
	make tidy
	make hugo

serve:
	hugo server -w --ignoreCache --logLevel info

serve-a:
	hugo server -w --ignoreCache --logLevel info -p 6060

serve-b:
	hugo server -w --ignoreCache --logLevel info -p 7070