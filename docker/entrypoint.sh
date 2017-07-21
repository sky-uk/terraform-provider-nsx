#!/bin/bash
govendor sync
while /bin/true; do fresh -c /fresh.conf ; sleep 5; done
