from faker import Faker
import random
import time



def str_time_prop(start, end, time_format, prop):
    """Get a time at a proportion of a range of two formatted times.

    start and end should be strings specifying times formatted in the
    given format (strftime-style), giving an interval [start, end].
    prop specifies how a proportion of the interval to be taken after
    start.  The returned time will be in the specified format.
    """

    stime = time.mktime(time.strptime(start, time_format))
    etime = time.mktime(time.strptime(end, time_format))

    ptime = stime + prop * (etime - stime)

    return time.strftime(time_format, time.localtime(ptime))



class Data(object):
    def __init__(self):
        self.fake = Faker()
    
    def reference(self):
        return self.fake.bothify(text='#########')
    
    def bank_account(self):
        return self.fake.bothify(text='####-####-####-####')
    
    def amount(self):
        return random.randint(20, 9999999)

    def date(self):
        return 