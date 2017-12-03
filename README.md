Soos
====

Docker driven building, this is successor of https://github.com/ngParty/soos with simplified API.

## How it works:

1. Generate checksum from e.g. package.json file
2. Generate Dockerfile if missing
3. Check if build image is present on local
4. Pull image if missing
5. Check if build image is present on local
6. If still mising build image with backed in dependencies ( e.g. node_modules )
7. Run image with e.g. yarn build task
8. Push image to registry

## Features

1. Fast builds
2. Zero build downtime due to outages
2. Zero dependecy and small ( 495 KB )
