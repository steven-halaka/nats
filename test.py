import asyncio
import nats
import nats.js.api as jsapi
from nats.js.client import Msg

async def main():
    nc = await nats.connect("nats://a:a@localhost:4222")
    print(f"Connected to {nc.connected_url}")
    print(f"Connect urls: {nc._server_info['connect_urls']}") 
    js = nc.jetstream()

    stream = await js.add_stream(name="tk",
        config=jsapi.StreamConfig(
            name="tk",
            subjects=["*"],
            storage=jsapi.StorageType.FILE,
            num_replicas=3
        ),
    )

    # Subscribe
    async def message_handler(msg: Msg):
        print(f"Received a message on {msg.subject}: {msg.data.decode()}")
        await msg.ack()

    await js.subscribe("foo", 
                       durable="durable-foo",
                       cb=message_handler,
                       config=jsapi.ConsumerConfig(
                           durable_name="durable-foo",
                           ack_policy=jsapi.AckPolicy.EXPLICIT))

    # Publish
    for i in range(3):
        await js.publish("foo", f"Message {i + 1}".encode())

    await asyncio.sleep(1)
    print(stream.state)
    await nc.drain()

if __name__ == "__main__":
    asyncio.run(main())
