CREATE TABLE `fund`.`fund_record` (
  `ID`   bigint(20) NOT NULL AUTO_INCREMENT ,
  `FUND_CODE`  varchar(10) NOT NULL COMMENT '基金代码',
  `NAME`  varchar(100) NOT NULL COMMENT '基金名称' ,
  `OPERATE_TYPE` int NOT NULL COMMENT '操作类型(0:买;1:卖;10:转换买;11:转换卖;100:红利再投资;101:现金分红)',
  `AMOUNT` decimal(10,2) NOT NULL COMMENT '金额',
  `UNIT_PRICE` decimal(10,2) NOT NULL COMMENT '成交单价',
  `DATE` date NOT NULL COMMENT '操作日期',
  `QUANTITY` decimal(10,2) NOT NULL COMMENT '份额',
  `GAIN`   decimal(10,2)  NOT NULL COMMENT '涨幅',
  `PROFIT` decimal(10,2) NOT NULL COMMENT '盈利',
  `TOTAL_PROFIT` decimal(10,2) NOT NULL COMMENT '累积盈利',
  `TOTAL_PURCHASE_AMOUNT` decimal(10,2) NOT NULL COMMENT '累积投入金额',
  `CREATED_AT`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `UPDATED_AT`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `RECORD_VERSION`  int(11) NOT NULL DEFAULT '1',
  `IS_DELETED`   tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `IDX_FUND_CODE_JZRQ` (`FUND_CODE`,`DATE`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='基金操作记录';
