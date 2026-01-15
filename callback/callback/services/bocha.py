"""
专为 AI Agent 设计的多模态搜索工具集 (Bocha)

版本: 1.1
最后更新: 2026-01-14

此脚本将复杂的 Bocha AI Search 功能分解为一系列目标明确、参数极少的独立工具，
专为 AI Agent 调用而设计。Agent 只需根据任务意图（如常规搜索、查找结构化数据或时效性新闻）
选择合适的工具，无需理解复杂的参数组合。

核心特性:
- 强大多模态能力: 能同时返回网页、图片、AI总结、追问建议，以及丰富的“模态卡”结构化数据。
- 模态卡支持: 针对天气、股票、汇率、百科、医疗等特定查询，可直接返回结构化数据卡片，便于Agent直接解析和使用。

主要工具:
- comprehensive_search: 执行全面搜索，返回网页、图片、AI总结及可能的模态卡。
- search_for_structured_data: 专门用于查询天气、股票、汇率等可触发“模态卡”的结构化信息。
- web_search_only: 执行纯网页搜索，不请求AI总结，速度更快。
- search_last_24_hours: 获取过去24小时内的最新信息。
- search_last_week: 获取过去一周内的主要报道。
"""

import datetime
import http
import json
import os
import sys
from email import header
from typing import Any, Dict, List, Literal, Optional

from regex import B, P

current_dir = os.path.dirname(os.path.abspath(__file__))
root_dir = os.path.dirname(os.path.dirname(current_dir))
if root_dir not in sys.path:
    print("添加callback根目录到Python路径:", root_dir)
    sys.path.append(root_dir)

from utils.log import logger
from utils.response import BizError

# 运行前请确保已安装 requests 库: pip install requests
try:
    import requests
except ImportError:
    raise ImportError("requests 库未安装，请运行 `pip install requests` 进行安装。")


# --- 1. 数据结构定义 ---
from dataclasses import dataclass, field


@dataclass
class WebpageResult:
    """网页搜索结果"""

    name: str
    url: str
    snippet: str
    display_url: Optional[str] = None
    date_last_crawled: Optional[str] = None


@dataclass
class ImageResult:
    """图片搜索结果"""

    name: str
    content_url: str
    host_page_url: Optional[str] = None
    thumbnail_url: Optional[str] = None
    width: Optional[int] = None
    height: Optional[int] = None


@dataclass
class ModalCardResult:
    """
    模态卡结构化数据结果
    这是 Bocha 搜索的核心特色，用于返回特定类型的结构化信息。
    """

    card_type: str  # 例如: weather_china, stock, baike_pro, medical_common
    content: Dict[str, Any]  # 解析后的JSON内容


@dataclass
class BochaResponse:
    """封装 Bocha API 的完整返回结果，以便在工具间传递"""

    query: str
    conversation_id: Optional[str] = None
    answer: Optional[str] = None  # AI生成的总结答案
    follow_ups: List[str] = field(default_factory=list)  # AI生成的追问
    webpages: List[WebpageResult] = field(default_factory=list)
    images: List[ImageResult] = field(default_factory=list)
    modal_cards: List[ModalCardResult] = field(default_factory=list)


# --- 2. 核心客户端与专用工具集 ---


class BochaMultimodalSearch:
    """
    一个包含多种专用多模态搜索工具的客户端。
    每个公共方法都设计为供 AI Agent 独立调用的工具。
    """

    BOCHA_BASE_URL = "https://api.bocha.cn/v1/ai-search"

    def __init__(self):
        """
        初始化客户端。
        """

    def _parse_search_response(
        self, response_dict: Dict[str, Any], query: str
    ) -> BochaResponse:
        """从API的原始字典响应中解析出结构化的BochaResponse对象"""

        final_response = BochaResponse(query=query)
        final_response.conversation_id = response_dict.get("conversation_id")

        messages = response_dict.get("messages", [])
        for msg in messages:
            role = msg.get("role")
            if role != "assistant":
                continue

            msg_type = msg.get("type")
            content_type = msg.get("content_type")
            content_str = msg.get("content", "{}")

            try:
                content_data = json.loads(content_str)
            except json.JSONDecodeError:
                # 如果内容不是合法的JSON字符串（例如纯文本的answer），则直接使用
                content_data = content_str

            if msg_type == "answer" and content_type == "text":
                final_response.answer = content_data

            elif msg_type == "follow_up" and content_type == "text":
                final_response.follow_ups.append(content_data)

            elif msg_type == "source":
                if content_type == "webpage":
                    web_results = content_data.get("value", [])
                    for item in web_results:
                        final_response.webpages.append(
                            WebpageResult(
                                name=item.get("name"),
                                url=item.get("url"),
                                snippet=item.get("snippet"),
                                display_url=item.get("displayUrl"),
                                date_last_crawled=item.get("dateLastCrawled"),
                            )
                        )
                elif content_type == "image":
                    final_response.images.append(
                        ImageResult(
                            name=content_data.get("name"),
                            content_url=content_data.get("contentUrl"),
                            host_page_url=content_data.get("hostPageUrl"),
                            thumbnail_url=content_data.get("thumbnailUrl"),
                            width=content_data.get("width"),
                            height=content_data.get("height"),
                        )
                    )
                # 所有其他 content_type 都视为模态卡
                else:
                    final_response.modal_cards.append(
                        ModalCardResult(card_type=content_type, content=content_data)
                    )

        return final_response

    def _search_internal(self, **kwargs) -> BochaResponse:
        """内部通用的搜索执行器，所有工具最终都调用此方法"""
        query = kwargs.get("query", "Unknown Query")
        header = kwargs.get("header", "Unknown Header")
        payload = {
            "stream": False,  # Agent工具通常使用非流式以获取完整结果
        }
        payload.update(kwargs)

        try:
            # logger.info(f"发送搜索请求，查询: {query} | header: {json.dumps(header, ensure_ascii=False)}")
            response = requests.post(
                self.BOCHA_BASE_URL, headers=header, json=payload, timeout=30
            )
            response.raise_for_status()  # 如果HTTP状态码是4xx或5xx，则抛出异常

            response_dict = response.json()

            return self._parse_search_response(response_dict, query)

        except requests.exceptions.RequestException as e:
            err = response.json()
            logger.exception(f"处理响应时发生未知错误: {err.get("code")}")
            raise BizError(err.get("message"), code=err.get("code"))
        except Exception as e:
            err = response.json()
            logger.exception(f"处理响应时发生未知错误: {err}")
            raise BizError(err.get("message"), code=err.get("code"))

    # --- Agent 可用的工具方法 ---

    def comprehensive_search(
        self, api_key: str, query: str, max_results: int = 10
    ) -> BochaResponse:
        """
        【工具】全面综合搜索: 执行一次标准的、包含所有信息类型的综合搜索。
        返回网页、图片、AI总结、追问建议和可能的模态卡。这是最常用的通用搜索工具。
        Agent可提供搜索查询(query)和可选的最大结果数(max_results)。
        """
        logger.info(f"--- TOOL: 全面综合搜索 (query: {query}) ---")
        return self._search_internal(
            header=get_bocha_header(api_key),
            query=query,
            count=max_results,
            answer=True,  # 开启AI总结
        )

    def web_search_only(
        self, api_key: str, query: str, max_results: int = 15
    ) -> BochaResponse:
        """
        【工具】纯网页搜索: 只获取网页链接和摘要，不请求AI生成答案。
        适用于需要快速获取原始网页信息，而不需要AI额外分析的场景。速度更快，成本更低。
        """
        logger.info(f"--- TOOL: 纯网页搜索 (query: {query}) ---")
        return self._search_internal(
            header=get_bocha_header(api_key),
            query=query,
            count=max_results,
            answer=False,  # 关闭AI总结
        )

    def search_for_structured_data(self, api_key: str, query: str) -> BochaResponse:
        """
        【工具】结构化数据查询: 专门用于可能触发“模态卡”的查询。
        当Agent意图是查询天气、股票、汇率、百科定义、火车票、汽车参数等结构化信息时，应优先使用此工具。
        它会返回所有信息，但Agent应重点关注结果中的 `modal_cards` 部分。
        """
        logger.info(f"--- TOOL: 结构化数据查询 (query: {query}) ---")
        # 实现上与 comprehensive_search 相同，但通过命名和文档引导Agent的意图
        return self._search_internal(
            header=get_bocha_header(api_key),
            query=query,
            count=5,  # 结构化查询通常不需要太多网页结果
            answer=True,
        )

    def search_last_24_hours(self, api_key: str, query: str) -> BochaResponse:
        """
        【工具】搜索24小时内信息: 获取关于某个主题的最新动态。
        此工具专门查找过去24小时内发布的内容。适用于追踪突发事件或最新进展。
        """
        logger.info(f"--- TOOL: 搜索24小时内信息 (query: {query}) ---")
        return self._search_internal(
            header=get_bocha_header(api_key),
            query=query,
            freshness="oneDay",
            answer=True,
        )

    def search_last_week(self, api_key: str, query: str) -> BochaResponse:
        """
        【工具】搜索本周信息: 获取关于某个主题过去一周内的主要报道。
        适用于进行周度舆情总结或回顾。
        """
        logger.info(f"--- TOOL: 搜索本周信息 (query: {query}) ---")
        return self._search_internal(
            header=get_bocha_header(api_key),
            query=query,
            freshness="oneWeek",
            answer=True,
        )


def get_bocha_header(api_key: str) -> Dict[str, str]:
    """生成Bocha API请求头"""
    return {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
        "Accept": "*/*",
    }


def print_response_summary(response):
    """简化的打印函数，用于展示测试结果"""
    if not response or not response.query:
        logger.error("未能获取有效响应。")
        return

    logger.info(f"\n查询: '{response.query}' | 会话ID: {response.conversation_id}")
    if hasattr(response, "answer") and response.answer:
        logger.info(f"AI摘要: {response.answer[:150]}...")

    logger.info(f"找到 {len(response.webpages)} 个网页")
    if hasattr(response, "images"):
        logger.info(f"找到 {len(response.images)} 张图片")
    if hasattr(response, "modal_cards"):
        logger.info(f"找到 {len(response.modal_cards)} 个模态卡")

    if hasattr(response, "modal_cards") and response.modal_cards:
        first_card = response.modal_cards[0]
        logger.info(f"第一个模态卡类型: {first_card.card_type}")

    if response.webpages:
        first_result = response.webpages[0]
        logger.info(f"第一条网页结果: {first_result.name}")

    if hasattr(response, "follow_ups") and response.follow_ups:
        logger.info(f"建议追问: {response.follow_ups}")

    logger.info("-" * 60)


if __name__ == "__main__":
    # 在运行前，请确保您已设置 BOCHA_API_KEY 环境变量

    try:
        # 初始化多模态搜索客户端，它内部包含了所有工具
        search_client = BochaMultimodalSearch()

        api_key = "sk-test-XXXXXXXXXXXXXXXXXXXXXX"  # 请替换为您的Bocha API Key

        # 场景1: Agent进行一次常规的、需要AI总结的综合搜索
        response1 = search_client.comprehensive_search(
            api_key=api_key, query="人工智能对未来教育的影响"
        )
        print_response_summary(response1)

        # 场景2: Agent需要查询特定结构化信息 - 天气
        if isinstance(search_client, BochaMultimodalSearch):
            response2 = search_client.search_for_structured_data(
                api_key=api_key, query="上海明天天气怎么样"
            )
            print_response_summary(response2)
            # 深度解析第一个模态卡
            if (
                response2.modal_cards
                and response2.modal_cards[0].card_type == "weather_china"
            ):
                logger.info(
                    "天气模态卡详情:",
                    json.dumps(
                        response2.modal_cards[0].content, indent=2, ensure_ascii=False
                    ),
                )

        # 场景3: Agent需要查询特定结构化信息 - 股票
        if isinstance(search_client, BochaMultimodalSearch):
            response3 = search_client.search_for_structured_data(
                api_key=api_key, query="东方财富股票"
            )
            print_response_summary(response3)

        # 场景4: Agent需要追踪某个事件的最新进展
        response4 = search_client.search_last_24_hours(
            api_key=api_key, query="C929大飞机最新消息"
        )
        print_response_summary(response4)

        # 场景5: Agent只需要快速获取网页信息，不需要AI总结
        if isinstance(search_client, BochaMultimodalSearch):
            response5 = search_client.web_search_only(
                api_key=api_key, query="Python dataclasses用法"
            )
            print_response_summary(response5)

        # 场景6: Agent需要回顾一周内关于某项技术的新闻
        response6 = search_client.search_last_week(
            api_key=api_key, query="量子计算商业化"
        )
        print_response_summary(response6)
    except ValueError as e:
        logger.exception(f"初始化失败: {e}")
        logger.error("请确保 BOCHA_API_KEY 环境变量已正确设置。")
    except Exception as e:
        logger.exception(f"测试过程中发生未知错误: {e}")
