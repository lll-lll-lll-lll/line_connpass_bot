# gcloud auth loginでログイン

# ①アーティファクトリポジトリ作成
create_new_artifact_repo:
	gcloud artifacts repositories create line-con-repo --location=asia-northeast1 --repository-format=docker

# ②`create_new_artifact_repo` makeでアーティファクトリポ作成してから、`deploy`でgroud run作成
deploy:
	gcloud builds submit --config=cloudbuild.yaml \
    --substitutions=_REPO_NAME=line-con-repo


# 作成したリポジトリ確認
repo_list:
	gcloud source repos list  

# Cloud Run service in a local development environment
_local_dev:
	gcloud beta code dev
