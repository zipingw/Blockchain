```python
#TODO：在Fabric chaincode中规定链上存储信息的类型
# 存储对象1：物联网数据
# {hash, source, time, store_nodes}
# hash: 文件哈希后的结果，类型为string
# source: 设备id，目前可以定类型为string
# time：数据生成的时间，类型为float等浮点数（时间戳本质是浮点数，可以看看chaincode支持什么样的小数形式）
# store_nodes: 存储节点id的拼接，可以参考"xxx,xxx"这样的格式，以字符串存储

# 存储对象2：操作记录
# {user_id, hash, type}
# user_id: string类型，发起操作的用户id
# hash：用户操作内容的hash值
# type：string类型，取值可以为{"query", "add", "modify"}等

#TODO: 安装如下应用及对应的python包：
# 1. RabbitMQ：如果用docker的话可以考虑用
#   docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_VHOST=myvhost -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=root rabbitmq
#   对应的python包：pika，安装方式：pip install pika
# 2. ipfshttpclient，详见https://pypi.org/project/ipfshttpclient/
# 3. fabric-sdk-py，详见https://fabric-sdk-py.readthedocs.io/en/latest/tutorial.html#operate-chaincodes-with-fabric-network


import pika


# MQTT subscriber
class TinyDBServer:
    package_size = 1    # 可调整参数：打包数据大小
    cache = {}

    # 测试add和select的时候可以不用开MQTT
    # def __init__(self):
        # self.connect()

    def connect():
        # 连接MQTT，其中的user_name, password, ip都需要替换
        user_info = pika.PlainCredentials('user_name', 'password')
        self.connection = pika.BlockingConnection(pika.ConnectionParameters('ip', 5672, 'myvhost', user_info))
    
    
    def run(self):
        channel = self.connection.channel()
        channel.queue_declare(queue='send_data')
        channel.queue_declare(queue='query_data')
        channel.basic_consume(queue='send_data', auto_ack=True, on_message_callback=self.send_data)
        channel.basic_consume(queue='query_data', auto_ack=True, on_message_callback=self.query_data)
        channel.start_consuming()

    
    def send_data():
        request = body.strip().decode('utf-8').split(',')
        data_source = request[0]
        data_time = request[1]
        content = request[2]
        if data_source in self.cache.keys():
            self.cache[data_source].append(content)
        else:
            self.cache[data_source] = [content]
        if len(self.cache[data_source]) >= self.package_size:
            self.add(','.join(self.cache[data_source]), data_source, data_time)
            self.cache[data_source] = []
    

    def query_data():
        request = body.strip().decode('utf-8').split(',')
        data_source = request[0]
        hash = request[1]
        return self.query(hash, data_source)



    #TODO: 请完善具体实现
    # input: content->需要存储的文件内容，我们假设用户在发送前已经完成加密，类型为string
    # output：hash->content的哈希值hash，类型为string
    def add(content, source, time):
        # 使用ipfshttpclient进行IPFS的存等操作
        ipfs = IPFS.client()

        # 将加密后的文件存储到IPFS中，得到IPFS返回的hash值
        hash = ipfs.add(c)          

        # 获取IPFS存储节点信息
        store_nodes = ipfs.find_where(hash)

        # 向Fabric存储信息，下面的chaincode.add可以自己写函数替代
        res = chaincode.add(hash, source, time, store_nodes)        

        # 判断Fabric存储是否成功（这里的逻辑需要你根据实际编写）
        if res is success:
            # 向Fabric记录操作
            chaincode.register(source, hash, "add")
            return hash
        else:
            return None


    #TODO: 请完善具体实现
    # input: hash->根据哈希值查找内容，hash类型为stirng
    # output：content->解密后的文件内容，类型为string
    def query(hash, source):
        # 使用ipfshttpclient进行IPFS的存等操作
        ipfs = IPFS.client()

        # 获取IPFS存储节点信息
        store_nodes = ipfs.find_where(hash)
        
        # 根据hash值向IPFS获取加密后的文件内容
        c = ipfs.cat(hash)

        # 验证链上存储的信息是否被IPFS修改，下面的chaincode.verify可以自己写函数替代
        flag = chaincode.verify(hash, c, store_nodes)

        # 向Fabric记录操作
        chaincode.register(source, hash, "query")

        # 如果没有被IPFS修改，则返回解密后的数据
        if flag is True:
            # AES解密，该encode函数可以从网上找AES已实现的包
            content = decode(c, key)
            return content
        else:
            return None

```



```python
from hfc.fabric import Client

cli = Client(net_profile="test/fixtures/network.json")

org1_admin = cli.get_user(org_name='org1.example.com', name='Admin') # get the admin user from local path

# 运行tinydb_server去访问 fabric network时 需要相应的 Cryptogen or Fabric-CA
# 从ca服务中获得 credential

from hfc.fabric_ca.caservice import ca_service

casvc = ca_service(target="http://127.0.0.1:7054")
# make local have the admin enrollment
adminEnrollment = casvc.enroll("admin", "adminpw")
# register a user to ca
secret = adminEnrollment.register("user1")
# make local have the user enrollment
user1Enrollment = casvc.enroll("user1", secret)
# now local will have the user reenrolled object 
user1ReEnrollment = casvc.reenroll(user1Enrollment)
RevokedCerts, CRL = adminEnrollment.revoke("user1")

```





```python
import ipfshttpclient
client = ipfshttpclient.connect()
res = client.add('test.txt')

with ipfshttpclient.connect() as client:
    hash = client.add('test.txt')['Hash']
    print(client.stat(hash))
    
```

