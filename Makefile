.PHONY: up

up:
	git add .
	git commit -am "update"
	git pull origin master
	git push origin master
	@echo "\n 代码提交发布..."

tag:
	git pull origin master
	git add .
	git commit -am "简单模板版本，支持windows"
	git push origin master
	git tag v1.0.0
	git push --tags
	@echo "\n tags 发布中..."