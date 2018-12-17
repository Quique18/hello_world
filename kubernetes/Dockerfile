FROM golang:latest

# Set the working directory to /app
#WORKDIR /app

# Copy the current directory contents into the container at /app
COPY rest /usr/bin/rest

# Install any needed packages specified in requirements.txt
# RUN pip install --trusted-host pypi.python.org -r requirements.txt

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV NAME Rest

# Run app.py when the container launches
CMD ["rest"]
