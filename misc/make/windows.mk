
misc/windows/oem/ssh.pub:
	echo "Copy your public SSH key here." > misc/windows/oem/ssh.pub

run-windows: misc/windows/oem/ssh.pub
	docker run -it --rm \
		--name amatl-windows \
		-p 8006:8006 \
		-p 2424:22 \
		--device=/dev/kvm \
		--device=/dev/net/tun \
		--cap-add NET_ADMIN \
		-v "$(CURDIR)/windows:/storage" \
		-v "$(CURDIR):/data/amatl" \
		-v "$(CURDIR)/misc/windows/oem:/oem" \
		-e LANGUAGE=French -e REGION=fr-FR -e KEYBOARD=fr-FR \
		--stop-timeout 120 \
		dockurr/windows

ssh-windows:
	ssh docker@127.0.0.1 -p 2424

reset-windows:
	sudo rm -rf "$(CURDIR)/windows"