FROM golang:1.18 AS BUILD_IMAGE

# Setting the Workspace
WORKDIR /usr/src/app

# Copying files to workspace
COPY . .

# Creating dist files to save build archives
RUN mkdir dist 

# Installing project depencencies
RUN go mod download

# Building
WORKDIR /usr/src/app/src
RUN go build -o ../dist

##############################################################################
FROM golang:1.18


# Copying build files to workspace
WORKDIR /usr/app/
COPY --from=BUILD_IMAGE /usr/src/app/dist ./

LABEL application=myapp

# Running the application
CMD ["./src"]