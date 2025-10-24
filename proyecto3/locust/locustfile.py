from locust import HttpUser, TaskSet, task, between
import random
import json

class MyTasks(TaskSet):
    
    @task(1)
    def engineering(self):
        # List of random municipalities
        municipalities = ["mixco", "guatemala", "amatitlan", "chinautla"]

        weather = ["sunny", "cloudy", "rainy", "foggy"]
    
        # Student data
        weather_data = {
            "municipality": random.choice(municipalities),  # Random municipality
            "temperature": random.randint(18, 28),  # Random temperature between 18 and 28
            "humidity": random.randint(40, 80),  # Random humidity between 40 and 80
            "weather": random.choice(weather)  # Random weather condition
        }
        
        # Send JSON as POST to the /engineering route
        headers = {'Content-Type': 'application/json'}
        self.client.post("/clima", json=weather_data, headers=headers)

    

class WebsiteUser(HttpUser):
    tasks = [MyTasks]
    wait_time = between(1, 5)  # Wait time between tasks (1 to 5 seconds)