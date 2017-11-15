all:
	mkdir -p output/bin output/conf
	cp -r conf/* output/conf/
	cp script/bootstrap.sh script/settings.py output
	go build -o output/bin/niwho_experiment_hellox

clean:
	rm -rf output

