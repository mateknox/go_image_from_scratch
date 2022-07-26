# Go docker image from scratch
Simple go app with static content, that can be build into a "scratch" image.

# Usage
**1)** docker build . -t go_image:1.0.0

**2)** docker run -d -p 5555:5555 <IMAGE_ID>

**3)** Go to 127.0.0.1:5555 to check if app is running.