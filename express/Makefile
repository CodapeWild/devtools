default: express cond

express:
	@goyacc -o comp_y.go -p comp comp.y

cond:
	@goyacc -o cond_y.go -p cond cond.y

clean:
	@rm -rf y.output
	@rm -rf express_y.go
	@rm -rf cond_y.go