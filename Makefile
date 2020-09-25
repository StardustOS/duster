image:
	sudo docker build -t duster-build .

duster:
	sudo docker run -v $PWD:/build:Z -it duster-build
