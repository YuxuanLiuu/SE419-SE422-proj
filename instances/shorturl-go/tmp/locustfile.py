from locust import HttpLocust, TaskSet, task
import string
import random

class UserTasks(TaskSet):
    def on_start(self):
        for _ in range(26):
            self.client.post("/",{"url":"sjtu.edu.cn"})

    @task(3)
    def insert(self):
        self.client.post("/",{"url":"sjtu.edu.cn"})
    @task(7)
    def index(self):
        self.client.get("/" + random.choice(string.ascii_letters.upper()) + "AAAAA")

class WebsiteUser(HttpLocust):
    task_set = UserTasks

