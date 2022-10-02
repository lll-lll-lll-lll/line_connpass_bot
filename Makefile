deploy:
	gcloud builds submit --config=cloudbuild.yaml \
    --substitutions=_REPO_NAME=line-con-repo


_create_new_artifact_repo:
	gcloud artifacts repositories create line-con-repo --location=asia-northeast1 --repository-format=docker

_repo_list:
	gcloud source repos list  