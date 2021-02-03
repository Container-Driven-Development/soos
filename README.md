Soos
====

Docker driven building, this is successor of https://github.com/ngParty/soos with simplified API.

## Why another strange named thingy? Why?

Because npm install is slow and I mean terribly slow when you try to use it properly with Docker for development.

Which means that you should run npm install every time you will create Docker container, or change branch, to have truly clean environment. When you have Docker images with node_modules baked in, you can skip this step and avoid all npm install issues. Only tricky part is how to build and run these images and how to mark them, so you, fellow developers and ci servers can reuse them.

Using Soos will simplify all of that, apart from building image it will also allow you to run image with right node_modules with single command in miliseconds. As a bonus you are forced into Docker driven development, which is speeding up onboarding and preventing issues caused by different environment setup.

## What it does then?

When you will run `soos` it will: 

1. Generate `checksum` from e.g. package.json file
2. Generate Dockerfile if missing
3. Check if build image is present on local
4. Pull image with `checksum` as a tag if exists
6. If image with `checksum` doesn't exist ( on first execution ) build image with backed in dependencies ( e.g. node_modules )
7. Run image
8. Push image to registry ( usually enabled on CI )

## Features

1. Fast builds
2. Zero build downtime due to newtork outages
2. Zero dependecy and small ( 495 KB )
