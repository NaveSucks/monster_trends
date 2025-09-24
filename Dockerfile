# Use a lightweight base image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary
COPY monster-trends .

# Copy templates + static into image
COPY templates ./templates
COPY static ./static

# Expose the port your Go app listens on
EXPOSE 3000

# Command to run the application
CMD ["./monster-trends"]