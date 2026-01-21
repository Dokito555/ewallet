apply-deployments:
	kubectl apply -f ./ums_deployment.yaml
	kubectl apply -f ./wallet_deployment.yaml
	kubectl apply -f ./transaction_deployment.yaml
	kubectl apply -f ./notification_deployment.yaml

apply-services:
	kubectl apply -f ./ums_service.yaml  
	kubectl apply -f ./wallet_service.yaml   
	kubectl apply -f ./transaction_service.yaml   
	kubectl apply -f ./notification_service.yaml    

delete-all:
	kubectl delete -f ./ums_deployment.yaml
	kubectl delete -f ./wallet_deployment.yaml
	kubectl delete -f ./transaction_deployment.yaml
	kubectl delete -f ./notification_deployment.yaml
	kubectl delete -f ./ums_service.yaml  
	kubectl delete -f ./wallet_service.yaml   
	kubectl delete -f ./transaction_service.yaml   
	kubectl delete -f ./notification_service.yaml 