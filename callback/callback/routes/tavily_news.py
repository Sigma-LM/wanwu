# app.py
import http
from dataclasses import asdict

from flask import Flask, g, jsonify, request

from callback.utils.decorators import require_api_key
from utils.log import logger
from utils.response import BizError

from . import callback_bp

# --- 路由接口 ---


@callback_bp.route("/tavily/news", methods=["POST"])
@require_api_key  # 直接使用导入的装饰器
def tavily_basic_search():
    """
    【工具】基础新闻搜索
    ---
    tags:
      - Tavily News
    summary: 执行基础新闻搜索
    description: 根据关键词搜索新闻，支持指定返回结果数量。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "OpenAI 最新动态"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 7)
                default: 7
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误 (缺少 query)
      401:
        description: API Key 无效或缺失
    """
    data = request.get_json() or {}
    query = data.get("query")

    if not query:
        raise BizError("Missing query", code=http.status_code.BAD_REQUEST)

    # g.agency 已经在装饰器里准备好了
    result = g.agency.basic_search_news(
        query=query, max_results=data.get("max_results")
    )
    logger.info(f"Tavily News Basic Search Result: {result}")
    return jsonify(result)


@callback_bp.route("/tavily/news/deep", methods=["POST"])
@require_api_key
def tavily_deep_search():
    """
    【工具】深度搜索
    ---
    tags:
      - Tavily News
    summary: 执行深度新闻搜索
    description: 对关键词进行更深入的挖掘和搜索。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "量子计算突破"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 20)
                default: 20
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误
    """
    data = request.get_json() or {}
    query = data.get("query")
    max_results = data.get("max_results")

    if not query:
        raise BizError("Missing 'query'", code=http.status_code.BAD_REQUEST)

    result = g.agency.deep_search_news(query=query, max_results=max_results)
    return jsonify(result)


@callback_bp.route("/tavily/news/day", methods=["POST"])
@require_api_key
def tavily_day_search():
    """
    【工具】搜索24小时内新闻
    ---
    tags:
      - Tavily News
    summary: 搜索最近一天的新闻
    description: 仅返回过去 24 小时内发布的相关新闻。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "纳斯达克今日行情"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 10)
                default: 10
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误
    """
    data = request.get_json() or {}
    query = data.get("query")
    max_results = data.get("max_results")

    if not query:
        raise BizError("Missing query", code=http.status_code.BAD_REQUEST)

    result = g.agency.search_news_last_24_hours(query=query, max_results=max_results)
    return jsonify(result)


@callback_bp.route("/tavily/news/week", methods=["POST"])
@require_api_key
def tavily_week_search():
    """
    【工具】搜索一周内新闻
    ---
    tags:
      - Tavily News
    summary: 搜索最近一周的新闻
    description: 仅返回过去一周内发布的相关新闻。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "纳斯达克今日行情"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 10)
                default: 10
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误
    """
    data = request.get_json() or {}
    query = data.get("query")
    max_results = data.get("max_results")

    if not query:
        raise BizError("Missing query", code=http.status_code.BAD_REQUEST)

    result = g.agency.search_news_last_week(query=query, max_results=max_results)
    return jsonify(result)


@callback_bp.route("/tavily/news/image", methods=["POST"])
@require_api_key
def tavily_image_search():
    """
    【工具】查找新闻图片
    ---
    tags:
      - Tavily News
    summary: 查找新闻图片
    description: 搜索与某个新闻主题相关的图片。此工具会返回图片链接及描述，适用于需要为报告或文章配图的场景。Agent只需提供搜索查询(query)。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "OpenAI 最新动态"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 5)
                default: 5
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误 (缺少 query)
      401:
        description: API Key 无效或缺失
    """
    data = request.get_json() or {}
    query = data.get("query")
    max_results = data.get("max_results")

    if not query:
        raise BizError("Missing query", code=http.status_code.BAD_REQUEST)

    result = g.agency.search_images_for_news(query=query, max_results=max_results)
    return jsonify(result)


@callback_bp.route("/tavily/news/date", methods=["POST"])
@require_api_key
def tavily_date_search():
    """
    【工具】按指定日期范围搜索新闻
    ---
    tags:
      - Tavily News
    summary: 按指定日期范围搜索新闻
    description: 在一个明确的历史时间段内搜索新闻。这是唯一需要Agent提供详细时间参数的工具。适用于需要对特定历史事件进行分析的场景。Agent需要提供查询(query)、开始日期(start_date)和结束日期(end_date)，格式均为 'YYYY-MM-DD'。
    parameters:
      - name: Authorization
        in: header
        description: "API Key"
        required: true
        schema:
          type: string
          default: "Bearer "
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - query
            properties:
              query:
                type: string
                description: 搜索关键词
                example: "OpenAI 最新动态"
              start_date:
                type: string
                description: 开始日期，格式 'YYYY-MM-DD'
                example: "2026-01-01"
              end_date:
                type: string
                description: 结束日期，格式 'YYYY-MM-DD'
                example: "2026-01-07"
              max_results:
                type: integer
                description: 返回结果的最大数量 (默认 15)
                default: 15
                example: 5
    responses:
      200:
        description: 搜索成功
        content:
          application/json:
            schema:
              type: object
              properties:
                query:
                  type: string
                  description: 搜索关键词
                answer:
                  type: string
                response_time:
                  type: number
                  format: float
                images:
                  type: array
                  items:
                    type: object
                    properties:
                      description:
                        type: string
                      url:
                        type: string
                results:
                  type: array
                  items:
                    type: object
                    properties:
                      title:
                        type: string
                      url:
                        type: string
                      content:
                        type: string
                      score:
                        type: number
                        format: float
                      raw_content:
                        type: string
                      published_date:
                        type: string
      400:
        description: 参数错误 (缺少 query)
      401:
        description: API Key 无效或缺失
    """
    data = request.get_json() or {}
    query = data.get("query")
    max_results = data.get("max_results")

    if not query:
        raise BizError("Missing query", code=http.status_code.BAD_REQUEST)

    result = g.agency.search_news_by_date(query=query, max_results=max_results)
    return jsonify(result)
