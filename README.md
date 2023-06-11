# AWS-SDK-GO





## Overview


This project provides `S3` file's `upload`, `download`,` delete`,` changeObjectStorage class`,  `getPresignedUrl` functionality using `aws-sdk-go` and integration with MongoDB for managing the related file metadata. Users can upload files and retrieve metadata for the uploaded files. 





## Getting Started


Follow the steps below to install and run the project





## Prerequisites


- Go language should be installed. (I'm using go1.20)
- You will need an AWS S3 bucket along with access key and secret key.
- You will need a MongoDB database connection URL and authentication credentials.





## Installation

1. Clone this project.


```shell
git clone https://github.com/seungbohong/aws-sdk-go.git
```

2. Navigate to the project directory.


```shell
cd your-project
```

3. Set up the environment variables.


   Create a `.env` file and enter the following information:

```plaintext
// AWS Credentials
AWS_ACCESS_KEY_ID=<Your AWS Access Key>
AWS_SECRET_ACCESS_KEY=<Your AWS Secret Access Key>

// MongoDB Credentials
MONGODB_USERNAME=<Your MongoDB Username>
MONGODB_PASSWORD=<Your MongoDB Password>
```

4. Install the required Go packages.


```shell
go mod tidy
```





## Running


1. Run the server.


```shell
./run.sh
```

2. Open your web browser and visit `http://localhost:8080`





## License


This project is licensed under the [MIT License](https://chat.openai.com/c/LICENSE).

