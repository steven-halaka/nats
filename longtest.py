import asyncio
import nats
import nats.js.api as jsapi
from datetime import datetime
from nats.js.client import Msg
from nats.aio.client import Client as NATS

async def disconnected_cb():
    print("[!] Disconnected from NATS!")

async def reconnected_cb():
    print("[+] Reconnected to NATS!")

async def closed_cb():
    print("[x] Connection closed.")

async def error_cb(e):
    print(f"[!] Error: {e}")

async def periodic_publisher(js):
    count = 1
    while True:
        try:
            msg = f"[{datetime.now().isoformat()}] Ping {count}"
            await js.publish("foo", msg.encode())
            print(f"Published: {msg}")
            count += 1
        except:
            print("publish error")
        await asyncio.sleep(1)

async def main():
    nc = NATS()
    await nc.connect("nats://a:a@localhost:4222",
                          disconnected_cb=disconnected_cb,
                          reconnected_cb=reconnected_cb,
                          closed_cb=closed_cb,
                          error_cb=error_cb)
    print(f"Connected to {nc.connected_url}")
    servers = [s.netloc for s in nc.servers]
    print(f"Connect urls: {servers}") 
    js = nc.jetstream()

    stream = await js.add_stream(name="tk",
        config=jsapi.StreamConfig(
            name="tk",
            subjects=["*"],
            storage=jsapi.StorageType.FILE,
            num_replicas=3
        ),
    )
    print(stream.state)

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

    asyncio.create_task(periodic_publisher(js))

    await asyncio.Event().wait()
    await nc.drain()

if __name__ == "__main__":
    asyncio.run(main())
