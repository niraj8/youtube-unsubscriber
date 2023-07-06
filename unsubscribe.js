// To nuke all your Youtube channels subscriptions
// Go to https://www.youtube.com/feed/channels
// Open Developer Console
// Copy and Paste this file's contents and Press Enter
// Refresh the page and repeat the process until all subscriptions are removed

async function scrollToBottom() {
    window.scrollTo(0, document.documentElement.scrollHeight);
    await delay(1000);

    if (window.innerHeight + window.scrollY < document.documentElement.scrollHeight) {
        await scrollToBottom();
    }
}

function delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function unsubscribe(item) {
    await item.click();

    // const unsubscribeXPath = '//*[@id="items"]/ytd-menu-service-item-renderer[4]/tp-yt-paper-item'
    const unsubscribeXPath = '//*[contains(text(), "Unsubscribe")]'
    var unsubscribeSnapshot = document.evaluate(unsubscribeXPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);

    if (unsubscribeSnapshot.singleNodeValue)
        unsubscribeSnapshot.singleNodeValue.click();
    else
        console.log("Element not found.");


    const confirmXPath = '//*[@id="confirm-button"]/yt-button-shape/button/yt-touch-feedback-shape/div/div[2]'
    const messageXPath = '//*[@id="scrollable"]/yt-formatted-string'

    var result = document.evaluate(confirmXPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);
    if (result.singleNodeValue) {
        result.singleNodeValue.click();
        console.log(document.evaluate(messageXPath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.textContent, "Yes")
    } else
        console.log("Element not found.");
}

async function unsubscribeAll(snapshot) {
    for (let i = 0; i <= snapshot.snapshotLength; i++) {
        await unsubscribe(snapshot.snapshotItem(i));
        await delay(1000);
    }
}

try {
    await scrollToBottom();
    console.log("Scrolling complete!");
} catch (error) {
    console.log("Error", error)
}

var snapshot = document.evaluate('//*[@id="notification-preference-button"]/ytd-subscription-notification-toggle-button-renderer-next/yt-button-shape/button', document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
await unsubscribeAll(snapshot)
