
run-windows:
	docker run -it --rm \
		--name amatl-windows \
		-p 8006:8006 --device=/dev/kvm \
		--device=/dev/net/tun \
		--cap-add NET_ADMIN \
		-v "$(CURDIR)/windows:/storage" \
		-v "$(CURDIR):/data/amatl" \
		-e LANGUAGE=French -e REGION=fr-FR -e KEYBOARD=fr-FR \
		--stop-timeout 120 \
		dockurr/windows

reset-windows:
	sudo rm -rf "$(CURDIR)/windows"