.PHONY: test

test:
	python3 -c "import yaml,os;v=yaml.safe_load(open('.gitlab-ci.yml', 'r').read()); [os.system(x) for x in (v['codecheck']['script'])]"
