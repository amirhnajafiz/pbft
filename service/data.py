from faker import Faker
import random
import time



def str_time_prop(start, end, time_format, prop):
    stime = time.mktime(time.strptime(start, time_format))
    etime = time.mktime(time.strptime(end, time_format))

    ptime = stime + prop * (etime - stime)

    return time.strftime(time_format, time.localtime(ptime))


def random_date(start, end, prop):
    return str_time_prop(start, end, '%m/%d/%Y %I:%M %p', prop)


class Data(object):
    def __init__(self):
        self.fake = Faker()
        self.banks = [
            'world bank',
            'bank of america',
            'blue bank',
            'chase bank',
            'visa card',
            'paypal',
            'american express',
            'bancrop',
            'saman',
            'first boston',
            'lehman brothers'
        ]
        self.types = ['atm', 'bank', 'app']
    
    def reference(self):
        return self.fake.bothify(text='#########')
    
    def bank_account(self):
        return self.fake.bothify(text='####-####-####-####')
    
    def amount(self):
        return random.randint(20, 9999999)

    def date(self):
        return random_date("1/1/2021 1:0 AM", "1/1/2022 11:59 PM", random.random())
    
    def bank(self):
        return random.choice(self.banks)
    
    def rtype(self):
        return random.choice(self.types)
