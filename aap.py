import requests

url = "https://www.example.com/api/"

headers = {"Content-Type": "application/json"}

response = requests.get(url, headers=headers)

print(response.json())