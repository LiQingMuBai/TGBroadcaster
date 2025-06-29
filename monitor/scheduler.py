import yaml
from telethon import TelegramClient
import asyncio
from datetime import datetime
import pytz
import logging

# 设置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


# 加载配置
def load_config():
    try:
        with open('config.yaml', 'r') as f:
            return yaml.safe_load(f)
    except Exception as e:
        logger.error(f"加载配置文件失败: {e}")
        raise


config = load_config()

# 初始化Telegram客户端
client = TelegramClient(
    'multi_group_scheduler',
    config['telegram']['api_id'],
    config['telegram']['api_hash']
)


async def send_to_group(group_id, message):
    """向指定群组发送消息"""
    try:
        await client.send_message(group_id, message)
        logger.info(f"消息成功发送到群组: {group_id}")
        return True
    except Exception as e:
        logger.error(f"发送到群组 {group_id} 失败: {e}")
        return False


async def send_scheduled_messages():
    """向所有群组发送定时消息"""
    message = config['schedule']['message']
    groups = config['groups']

    results = await asyncio.gather(
        *[send_to_group(group, message) for group in groups],
        return_exceptions=True
    )

    success_count = sum(1 for r in results if r is True)
    logger.info(f"消息发送完成: 成功 {success_count}/{len(groups)} 个群组")


async def scheduler():
    """定时任务调度器"""
    logger.info("定时消息服务已启动...")

    # 设置时区
    tz = pytz.timezone('Asia/Shanghai')

    while True:
        now = datetime.now(tz)
        current_time = now.strftime("%H:%M")

        if current_time in config['schedule']['times']:
            logger.info(f"触发定时发送: {current_time}")
            await send_scheduled_messages()

            # 发送后等待61秒，避免同一分钟内重复发送
            await asyncio.sleep(61)
        else:
            # 每分钟检查一次
            await asyncio.sleep(60 - now.second)


async def main():
    await client.start()
    me = await client.get_me()
    logger.info(f"登录成功，用户: {me.username} (ID: {me.id})")
    logger.info(f"监控的群组数量: {len(config['groups'])}")
    await scheduler()


if __name__ == '__main__':
    try:
        with client:
            client.loop.run_until_complete(main())
    except KeyboardInterrupt:
        logger.info("程序被用户中断")
    except Exception as e:
        logger.error(f"程序异常: {e}")