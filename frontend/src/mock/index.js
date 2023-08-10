

export function mock(params) {
    const tutorials = {
        list: [
            {
                id: 1,
                "repoUrl": "https://github.com/decert-me/blockchain-basic",
                "label": "区块链基础",
                "catalogueName": "blockchain-basic",
                "docType": "docusaurus",
                "img": "https://ipfs.decert.me/images/blockchain-basic.png",
                "desc": "区块链是一项令人兴奋且在快速发展的技术，你也许看到过这些频繁在社交媒体、新闻频道上冒出的新名词：智能合约、代币（通证）、Web3、DeFi、DAO 组织。 如果你还不是很明白他们的意思，这份免费区块链基础教程就是为你（小白们）准备的。",
                "challenge": 10000,
                "startPage": "blockchain-basic/start",
                "category": ["dapps"],
                "theme": ["defi"],
                "language": "zh",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 1
            },
            {
                id: 2,
                "repoUrl": "https://github.com/decert-me/learnsolidity",
                "label": "学习 Solidity",
                "catalogueName": "solidity",
                "img": "https://ipfs.decert.me/images/learn-solidity.png",
                "desc": "Solidity是一种专门为以太坊平台设计的编程语言，它是EVM智能合约的核心，是区块链开发人员必须掌握的一项技能。",
                "docType": "docusaurus",
                "challenge": 10002,
                "startPage": "solidity/intro",
                "category": ["chain-public"],
                "theme": ["btc"],
                "language": "zh",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 2,
            },
            {
                id: 3,
                "repoUrl": "https://github.com/SixdegreeLab/MasteringChainAnalytics",
                "label": "成为链上数据分析师",
                "catalogueName": "MasteringChainAnalytics",
                "img": "https://ipfs.learnblockchain.cn/images/sixdegree.png",
                "desc": "本教程是一个面向区块链爱好者的系列教程，帮助新手用户从零开始学习区块链数据分析，成为一名链上数据分析师。",
                "docType": "gitbook",
                "challenge": 10004,
                "commitHash": "cfb5403f14932b520adc9084673bd1c011f1aa2b",
                "startPage": "MasteringChainAnalytics/README",
                "category": ["layer2"],
                "theme": ["web3"],
                "language": "zh",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 3,
            },
            {
                id: 4,
                "repoUrl": "https://github.com/miguelmota/ethereum-development-with-go-book",
                "label": "用Go来做以太坊开发",
                "catalogueName": "ethereum-development-with-go-book",
                "img": "https://ipfs.learnblockchain.cn/images/ethereum-development-with-go.jpg",
                "desc": "这本迷你书的本意是给任何想用Go进行以太坊开发的同学一个概括的介绍。本意是如果你已经对以太坊和Go有一些熟悉，但是对于怎么把两者结合起来还有些无从下手，那这本书就是一个好的起点。",
                "docType": "gitbook",
                "branch": "master",
                "docPath": "/zh",
                "startPage": "ethereum-development-with-go-book/README",
                "category": ["safe"],
                "theme": ["blockchain"],
                "language": "zh",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 4,
            },
            {
                id: 5,
                "repoUrl": "https://github.com/RandyPen/sui-move-intro-course-zh",
                "label": "Sui Move 导论",
                "startPage": "sui-move-intro-course-zh/README",
                "catalogueName": "sui-move-intro-course-zh",
                "img": "https://ipfs.learnblockchain.cn/images/sui.png",
                "docType": "mdBook",
                "category": ["storage"],
                "theme": ["metaverse"],
                "language": "zh",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 5,
            },
            {
                id: 6,
                "repoUrl": "https://github.com/ingonyama-zk/ingopedia",
                "label": "ZKP Encyclopedia",
                "startPage": "ingopedia/communityguide",
                "catalogueName": "ingopedia",
                "img": "https://ipfs.learnblockchain.cn/images/helmet4a.png",
                "desc": "A curated list of ZK/FHE resources and links. ",
                "docType": "mdBook",
                "branch": "master",
                "docPath": "/src",
                "commitHash": "9f27ed7ab0fdd92c446a14e1df59891c17f6e8ed",
                "category": ["others"],
                "theme": ["cryptography", "eth", "dao"],
                "language": "en",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "0",
                key: 6,
            },
            {
                id: 7,
                "url": "https://www.youtube.com/watch?v=CTJ1JkYLiyw&list=PLBIxCe-LDnGC-hfuyKRBvi1eC_X_7JzU7",
                "label": "ZK Shanghai-Lesson2 零知识证明工作坊",
                "catalogueName": "ZK",
                "docType": "video",
                "img": "https://ipfs.decert.me/images/blockchain-basic.png",
                "desc": "区块链是一项令人兴奋且在快速发展的技术，你也许看到过这些频繁在社交媒体、新闻频道上冒出的新名词：智能合约、代币（通证）、Web3、DeFi、DAO 组织。 如果你还不是很明白他们的意思，这份免费区块链基础教程就是为你（小白们）准备的。",
                "challenge": 10000,
                "sort": "ordered",
                "videoCategory": "youtube",
                "category": ["others"],
                "theme": ["cryptography", "eth", "dao"],
                "language": "en",
                "time": 9000000,  //  预估时间
                "difficulty": 2,     //  难度
                status: "1",
                key: 7,
            }
        ]
    }

    return {
        tutorials
    }
}