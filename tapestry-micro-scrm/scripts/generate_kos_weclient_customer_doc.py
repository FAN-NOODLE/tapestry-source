#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""生成并校验《KOS 与 WeClient 素材数据同步技术说明》。"""

from __future__ import annotations

import argparse
import html
import xml.etree.ElementTree as ET
from pathlib import Path
from zipfile import ZIP_DEFLATED, ZipFile


ROOT = Path(__file__).resolve().parent.parent
DEFAULT_OUTPUT = ROOT / "docs" / "customer" / "KOS与WeClient素材数据同步技术说明.docx"


def escape(value: str) -> str:
    return html.escape(value, quote=False)


def run(text: str, *, bold: bool = False, color: str = "333333", size: int = 21) -> str:
    properties = [
        '<w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:eastAsia="Microsoft YaHei"/>',
        f'<w:color w:val="{color}"/>',
        f'<w:sz w:val="{size}"/><w:szCs w:val="{size}"/>',
    ]
    if bold:
        properties.append("<w:b/><w:bCs/>")
    return f'<w:r><w:rPr>{"".join(properties)}</w:rPr><w:t xml:space="preserve">{escape(text)}</w:t></w:r>'


def paragraph(
    text: str = "",
    *,
    style: str = "BodyText",
    align: str | None = None,
    bold: bool = False,
    color: str = "333333",
    size: int = 21,
    before: int = 0,
    after: int = 120,
) -> str:
    alignment = f'<w:jc w:val="{align}"/>' if align else ""
    properties = (
        f'<w:pStyle w:val="{style}"/><w:spacing w:before="{before}" w:after="{after}" '
        f'w:line="360" w:lineRule="auto"/>{alignment}'
    )
    return f'<w:p><w:pPr>{properties}</w:pPr>{run(text, bold=bold, color=color, size=size)}</w:p>'


def rich_paragraph(parts: list[tuple[str, bool, str]], *, align: str | None = None) -> str:
    alignment = f'<w:jc w:val="{align}"/>' if align else ""
    runs = "".join(run(text, bold=bold, color=color) for text, bold, color in parts)
    return (
        '<w:p><w:pPr><w:pStyle w:val="BodyText"/><w:spacing w:after="120" '
        f'w:line="360" w:lineRule="auto"/>{alignment}</w:pPr>{runs}</w:p>'
    )


def heading(text: str, *, level: int = 1) -> str:
    if level == 1:
        return paragraph(text, style="Heading1", bold=True, color="17365D", size=32, before=80, after=220)
    return paragraph(text, style="Heading2", bold=True, color="2F5597", size=25, before=80, after=120)


def bullet(text: str) -> str:
    return rich_paragraph([("•  ", True, "2F5597"), (text, False, "333333")])


def page_break() -> str:
    return '<w:p><w:r><w:br w:type="page"/></w:r></w:p>'


def cell(text: str, *, fill: str = "FFFFFF", bold: bool = False, color: str = "333333", width: int = 3000) -> str:
    content = paragraph(text, align="center", bold=bold, color=color, size=19, after=40)
    return (
        f'<w:tc><w:tcPr><w:tcW w:w="{width}" w:type="dxa"/><w:shd w:fill="{fill}"/>'
        '<w:tcMar><w:top w:w="90" w:type="dxa"/><w:left w:w="90" w:type="dxa"/>'
        '<w:bottom w:w="90" w:type="dxa"/><w:right w:w="90" w:type="dxa"/></w:tcMar>'
        f'</w:tcPr>{content}</w:tc>'
    )


def table(rows: list[list[tuple[str, str, bool, str]]], widths: list[int]) -> str:
    grid = "".join(f'<w:gridCol w:w="{width}"/>' for width in widths)
    rendered_rows = []
    for row in rows:
        cells = "".join(
            cell(text, fill=fill, bold=bold, color=color, width=widths[index])
            for index, (text, fill, bold, color) in enumerate(row)
        )
        rendered_rows.append(f'<w:tr>{cells}</w:tr>')
    borders = "".join(
        f'<w:{edge} w:val="single" w:sz="4" w:space="0" w:color="D9E2F3"/>'
        for edge in ("top", "left", "bottom", "right", "insideH", "insideV")
    )
    return (
        '<w:tbl><w:tblPr><w:tblW w:w="0" w:type="auto"/><w:jc w:val="center"/>'
        f'<w:tblBorders>{borders}</w:tblBorders><w:tblCellMar><w:top w:w="40" w:type="dxa"/>'
        '<w:left w:w="40" w:type="dxa"/><w:bottom w:w="40" w:type="dxa"/>'
        f'<w:right w:w="40" w:type="dxa"/></w:tblCellMar></w:tblPr><w:tblGrid>{grid}</w:tblGrid>'
        f'{"".join(rendered_rows)}</w:tbl>'
    )


def info_table() -> str:
    header = "4472C4"
    label = "D9EAF7"
    rows = [
        [("项目", header, True, "FFFFFF"), ("说明", header, True, "FFFFFF")],
        [("请求方式", label, True, "17365D"), ("GET", "FFFFFF", False, "333333")],
        [("接口路径", label, True, "17365D"), ("/api/biz/coach_mark_post_results", "FFFFFF", False, "333333")],
        [("查询参数", label, True, "17365D"), ("updated_at、page、size", "FFFFFF", False, "333333")],
        [("鉴权请求头", label, True, "17365D"), ("X-Coach-Mark-Token", "FFFFFF", False, "333333")],
        [("成功业务码", label, True, "17365D"), ("code = 0", "FFFFFF", False, "333333")],
        [("数据列表", label, True, "17365D"), ("data.list", "FFFFFF", False, "333333")],
    ]
    return table(rows, [2500, 6500])


def sequence_table() -> str:
    empty = ("│", "F7F9FC", False, "A6A6A6")
    rows = [
        [("KOS", "44546A", True, "FFFFFF"), ("WeClient 定时任务", "4472C4", True, "FFFFFF"), ("WeClient 素材库", "44546A", True, "FFFFFF")],
        [("创建素材", "E2F0D9", True, "385723"), empty, empty],
        [("返回分页数据  ───────────────→", "DDEBF7", False, "2F5597"), ("请求 7 天窗口内数据  ←────", "DDEBF7", False, "2F5597"), empty],
        [empty, ("按 post_id 查询  ─────────→", "DDEBF7", False, "2F5597"), ("未找到", "F7F9FC", False, "666666")],
        [empty, ("新增素材  ─────────────────→", "E2F0D9", True, "385723"), ("新增完成", "E2F0D9", True, "385723")],
        [("更新同一素材", "FCE4D6", True, "9C5700"), empty, empty],
        [("返回更新数据  ───────────────→", "DDEBF7", False, "2F5597"), ("下一次定时请求  ←────────", "DDEBF7", False, "2F5597"), empty],
        [empty, ("按 post_id 查询  ─────────→", "DDEBF7", False, "2F5597"), ("已存在", "F7F9FC", False, "666666")],
        [empty, ("更新素材  ─────────────────→", "FCE4D6", True, "9C5700"), ("更新完成；N 时禁用", "FCE4D6", True, "9C5700")],
        [empty, ("当前页完成后继续下一页", "FFF2CC", False, "7F6000"), empty],
        [empty, ("失败记录异常；后续任务重新拉取", "FFF2CC", False, "7F6000"), empty],
    ]
    return table(rows, [3000, 3600, 3000])


def build_document_xml() -> str:
    content: list[str] = []

    # 第 1 页：封面
    content.extend(
        [
            paragraph("KOS × WeClient", align="center", bold=True, color="4472C4", size=28, before=1500, after=380),
            paragraph("KOS 与 WeClient 素材数据同步技术说明", align="center", bold=True, color="17365D", size=42, after=650),
            paragraph("面向客户的系统对接与数据同步说明", align="center", color="666666", size=24, after=500),
            table(
                [
                    [("文档版本", "D9EAF7", True, "17365D"), ("V1.0", "FFFFFF", False, "333333")],
                    [("发布日期", "D9EAF7", True, "17365D"), ("2026-07-16", "FFFFFF", False, "333333")],
                    [("适用系统", "D9EAF7", True, "17365D"), ("KOS、WeClient", "FFFFFF", False, "333333")],
                ],
                [3000, 5200],
            ),
            paragraph("本文档用于说明双方系统的素材增量同步机制。接口鉴权信息由双方通过安全渠道另行提供。", align="center", color="7F7F7F", size=18, before=520, after=0),
        ]
    )
    content.append(page_break())

    # 第 2 页：对接概览
    content.extend(
        [
            heading("1. 对接概览"),
            heading("1.1 文档目的", level=2),
            paragraph("本对接用于将 KOS 中创建或更新的素材同步至 WeClient。KOS 提供标准查询接口，WeClient 通过定时任务主动拉取并完成素材新增或更新。"),
            heading("1.2 双方职责", level=2),
            table(
                [
                    [("系统", "4472C4", True, "FFFFFF"), ("主要职责", "4472C4", True, "FFFFFF")],
                    [("KOS", "D9EAF7", True, "17365D"), ("维护素材数据；提供增量查询接口；返回素材有效状态及商品关联信息。", "FFFFFF", False, "333333")],
                    [("WeClient", "D9EAF7", True, "17365D"), ("按计划拉取数据；分页处理；按 post_id 新增或更新素材；记录并重试异常。", "FFFFFF", False, "333333")],
                ],
                [2300, 6700],
            ),
            heading("1.3 总体流程", level=2),
            bullet("KOS 创建或更新素材，并在接口中提供最新数据。"),
            bullet("WeClient 定时携带 updated_at、page、size 调用 KOS 接口。"),
            bullet("WeClient 按 post_id 判断素材是首次同步还是已有数据更新。"),
            bullet("同步完成后，素材进入 WeClient 对应月份分组，供后续业务使用。"),
            paragraph("同步方向：KOS  →  WeClient", align="center", bold=True, color="4472C4", size=26, before=240, after=0),
        ]
    )
    content.append(page_break())

    # 第 3 页：接口与增量规则
    content.extend(
        [
            heading("2. 接口调用与增量规则"),
            heading("2.1 接口概览", level=2),
            info_table(),
            heading("2.2 调用示例", level=2),
            paragraph("https://kosuat.coach.com.cn/api/api/biz/coach_mark_post_results", style="CodeBlock", color="1F4E79", size=18, after=40),
            paragraph("?updated_at=2026-07-05%2020:00:00&page=1&size=1000", style="CodeBlock", color="1F4E79", size=18, after=100),
            paragraph("请求头示例：X-Coach-Mark-Token: ********（实际值通过安全渠道提供）", style="CodeBlock", color="1F4E79", size=18),
            heading("2.3 增量与分页", level=2),
            bullet("WeClient 每次默认查询当前时间前 7 天以来发生更新的数据。"),
            bullet("接口按页返回数据，WeClient 完成当前页后继续请求下一页，直至本批次处理完毕。"),
            bullet("7 天窗口允许同一素材重复返回；WeClient 按 post_id 幂等处理，不会重复创建。"),
            bullet("返回成功时业务码为 code = 0，素材数据位于 data.list。"),
        ]
    )
    content.append(page_break())

    # 第 4 页：时序图
    content.extend(
        [
            heading("3. 素材同步时序"),
            paragraph("下图展示 KOS 创建素材后 WeClient 新增，以及 KOS 更新同一素材后 WeClient 更新的连续过程。", after=180),
            sequence_table(),
            paragraph("图 1  KOS 与 WeClient 素材同步时序", align="center", color="666666", size=18, before=120, after=80),
            paragraph("说明：重复返回的数据由 WeClient 按 post_id 幂等处理；当 is_valid = N 时，对应素材在 WeClient 中设置为禁用。", style="Note", color="7F6000", size=18, after=0),
        ]
    )
    content.append(page_break())

    # 第 5 页：业务规则与异常补偿
    content.extend(
        [
            heading("4. 业务规则与异常处理"),
            heading("4.1 新增与更新", level=2),
            bullet("post_id 尚未同步：WeClient 新增素材及其商品关联信息。"),
            bullet("post_id 已经存在：WeClient 使用 KOS 最新返回内容更新已有素材。"),
            bullet("is_valid = N：WeClient 将对应素材设置为禁用；其他值按有效素材处理。"),
            heading("4.2 分组、有效期与商品编号", level=2),
            bullet("素材按照 data_date 对应月份归入月份分组；分组不存在时由 WeClient 自动创建。"),
            bullet("同步素材的默认有效期截至 2035-12-31。"),
            bullet("product_id 用于建立商品关联；单个值包含 & 时，WeClient 将其拆分为多个商品编号。"),
            bullet("商品详情由 WeClient 在素材同步流程中补全；未查询到商品时忽略，不影响素材同步结果。"),
            heading("4.3 异常与补偿", level=2),
            bullet("接口调用或单条数据处理失败时，WeClient 记录异常信息，便于定位。"),
            bullet("后续定时任务继续使用 7 天重叠窗口，重新拉取并处理相关数据。"),
            bullet("KOS 保证接口可访问、鉴权信息有效，并在数据更新后返回最新记录。"),
            paragraph("对接结论：KOS 负责提供最新素材，WeClient 负责周期性拉取、幂等落库与失败补偿。", style="Note", bold=True, color="17365D", size=21, before=220, after=0),
        ]
    )

    section = (
        '<w:sectPr><w:footerReference w:type="default" r:id="rId3"/>'
        '<w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="1134" w:right="1134" '
        'w:bottom="1134" w:left="1134" w:header="600" w:footer="600" w:gutter="0"/>'
        '<w:cols w:space="720"/><w:docGrid w:linePitch="312"/></w:sectPr>'
    )
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" '
        'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">'
        f'<w:body>{"".join(content)}{section}</w:body></w:document>'
    )


def build_styles_xml() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:docDefaults><w:rPrDefault><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:eastAsia="Microsoft YaHei"/><w:sz w:val="21"/><w:szCs w:val="21"/></w:rPr></w:rPrDefault></w:docDefaults>
  <w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/><w:qFormat/></w:style>
  <w:style w:type="paragraph" w:styleId="BodyText"><w:name w:val="正文"/><w:basedOn w:val="Normal"/><w:qFormat/><w:pPr><w:spacing w:after="120" w:line="360" w:lineRule="auto"/></w:pPr></w:style>
  <w:style w:type="paragraph" w:styleId="Heading1"><w:name w:val="标题 1"/><w:basedOn w:val="Normal"/><w:qFormat/><w:pPr><w:keepNext/><w:outlineLvl w:val="0"/></w:pPr></w:style>
  <w:style w:type="paragraph" w:styleId="Heading2"><w:name w:val="标题 2"/><w:basedOn w:val="Normal"/><w:qFormat/><w:pPr><w:keepNext/><w:outlineLvl w:val="1"/></w:pPr></w:style>
  <w:style w:type="paragraph" w:styleId="CodeBlock"><w:name w:val="接口示例"/><w:basedOn w:val="Normal"/><w:pPr><w:shd w:fill="F2F2F2"/><w:ind w:left="180" w:right="180"/><w:spacing w:before="60" w:after="60"/></w:pPr></w:style>
  <w:style w:type="paragraph" w:styleId="Note"><w:name w:val="说明"/><w:basedOn w:val="Normal"/><w:pPr><w:shd w:fill="FFF2CC"/><w:ind w:left="180" w:right="180"/><w:spacing w:before="80" w:after="80"/></w:pPr></w:style>
</w:styles>"""


def build_content_types() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
  <Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>
  <Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
  <Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
</Types>"""


def build_package_relationships() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>"""


def build_document_relationships() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>
</Relationships>"""


def build_settings_xml() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:settings xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:zoom w:percent="100"/><w:defaultTabStop w:val="420"/><w:compat/>
</w:settings>"""


def build_footer_xml() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:color w:val="999999"/><w:sz w:val="18"/></w:rPr><w:t>KOS × WeClient  |  第 </w:t></w:r><w:fldSimple w:instr=" PAGE "><w:r><w:rPr><w:color w:val="999999"/><w:sz w:val="18"/></w:rPr><w:t>1</w:t></w:r></w:fldSimple><w:r><w:rPr><w:color w:val="999999"/><w:sz w:val="18"/></w:rPr><w:t> 页</w:t></w:r></w:p>
</w:ftr>"""


def build_core_properties() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <dc:title>KOS 与 WeClient 素材数据同步技术说明</dc:title><dc:creator>WeClient</dc:creator><dc:subject>素材数据同步</dc:subject><dc:description>KOS 与 WeClient 客户技术对接文档</dc:description><dcterms:created xsi:type="dcterms:W3CDTF">2026-07-16T00:00:00Z</dcterms:created><dcterms:modified xsi:type="dcterms:W3CDTF">2026-07-16T00:00:00Z</dcterms:modified>
</cp:coreProperties>"""


def build_app_properties() -> str:
    return """<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><Application>WeClient Documentation Generator</Application><Pages>5</Pages><Company>WeClient</Company></Properties>"""


def build_document(output_path: Path) -> None:
    output_path.parent.mkdir(parents=True, exist_ok=True)
    parts = {
        "[Content_Types].xml": build_content_types(),
        "_rels/.rels": build_package_relationships(),
        "word/document.xml": build_document_xml(),
        "word/styles.xml": build_styles_xml(),
        "word/settings.xml": build_settings_xml(),
        "word/footer1.xml": build_footer_xml(),
        "word/_rels/document.xml.rels": build_document_relationships(),
        "docProps/core.xml": build_core_properties(),
        "docProps/app.xml": build_app_properties(),
    }
    with ZipFile(output_path, "w", ZIP_DEFLATED) as archive:
        for name, content in parts.items():
            archive.writestr(name, content.encode("utf-8"))


def validate_document(output_path: Path) -> None:
    required_parts = {
        "[Content_Types].xml",
        "_rels/.rels",
        "word/document.xml",
        "word/styles.xml",
        "word/settings.xml",
    }
    required_text = {
        "KOS 与 WeClient 素材数据同步技术说明",
        "X-Coach-Mark-Token",
        "post_id",
        "7 天",
        "is_valid = N",
        "2035-12-31",
    }
    forbidden_text = {
        "REDACTED_FOR_SOURCE_SHARING",
        "source_link_kos",
        "source_relation_goods",
        "CREATE TABLE",
    }
    if not output_path.exists():
        raise ValueError(f"文档不存在：{output_path}")
    with ZipFile(output_path) as archive:
        broken = archive.testzip()
        if broken:
            raise ValueError(f"DOCX 压缩包损坏：{broken}")
        names = set(archive.namelist())
        missing_parts = required_parts - names
        if missing_parts:
            raise ValueError(f"缺少 OOXML 部件：{sorted(missing_parts)}")
        for name in names:
            if name.endswith(".xml"):
                ET.fromstring(archive.read(name))
        document = archive.read("word/document.xml").decode("utf-8")
    missing_text = {item for item in required_text if item not in document}
    if missing_text:
        raise ValueError(f"缺少关键文案：{sorted(missing_text)}")
    leaked_text = {item for item in forbidden_text if item in document}
    if leaked_text:
        raise ValueError(f"包含禁止公开内容：{sorted(leaked_text)}")
    if document.count('w:type="page"') != 4:
        raise ValueError("文档必须包含 4 个显式分页符，以形成 5 页内容")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--check", action="store_true", help="仅校验已生成文档")
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT, help="DOCX 输出路径")
    args = parser.parse_args()
    if not args.check:
        build_document(args.output)
    validate_document(args.output)
    print(f"文档校验成功：{args.output}")


if __name__ == "__main__":
    main()
