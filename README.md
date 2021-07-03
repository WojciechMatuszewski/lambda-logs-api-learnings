# Discovering Lambda Logs API

Me fiddling with the _AWS Lambda Logs API_ and discovering it's limitations.

## Deploying

This repository uses _AWS SAM_ for managing AWS-related resources.

To deploy the resources:

1. Make sure you have your AWS account set up and necessary environment variables present.

2. Run the deploy command

   ```sh
   make deploy
   ```

## Learnings

- You can test the **extension, NOT the Logs API** via `sam local invoke`

- Even though I'm not subscribing to any event from the _Extensions API_, I still had to make the request to the `event/next` route. Otherwise my extension would crash

- The **_Logs API_ is only available in the AWS environment. AFAIK there is no emulator available**

- The _extension loop_ is needed, otherwise you can miss the logs that are pushed to the endpoint you have specified
