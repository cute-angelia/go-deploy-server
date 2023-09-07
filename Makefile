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
	git commit -am "修改git日志获取方式；增加api更新拉取代码"
	git push origin master
	git tag v1.0.1
	git push --tags
	@echo "\n tags 发布中..."