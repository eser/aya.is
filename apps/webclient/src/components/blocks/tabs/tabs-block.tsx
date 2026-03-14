import React, { useState } from "react";
import styles from "./tabs-block.module.css";

interface TabBlockProps {
  label: string;
  children?: React.ReactNode;
}

function TabBlock(props: TabBlockProps) {
  return <>{props.children}</>;
}

interface TabsBlockProps {
  children?: React.ReactNode;
}

function TabsBlock(props: TabsBlockProps) {
  const [activeIndex, setActiveIndex] = useState(0);

  const childArray = React.Children.toArray(props.children);
  const tabs = childArray.filter(
    (child): child is React.ReactElement<TabBlockProps> =>
      React.isValidElement(child) && child.type === TabBlock,
  );

  if (tabs.length === 0) {
    return (
      <div className={styles.tabs}>
        <div className={styles.placeholder}>Add tabs</div>
      </div>
    );
  }

  const safeIndex = activeIndex < tabs.length ? activeIndex : 0;

  return (
    <div className={styles.tabs}>
      <div className={styles.tabList} role="tablist">
        {tabs.map((tab, index) => {
          const label =
            tab.props.label !== undefined &&
            tab.props.label !== null &&
            tab.props.label !== ""
              ? tab.props.label
              : `Tab ${index + 1}`;

          return (
            <button
              key={index}
              className={styles.tab}
              data-active={index === safeIndex}
              role="tab"
              aria-selected={index === safeIndex}
              onClick={() => setActiveIndex(index)}
              type="button"
            >
              {label}
            </button>
          );
        })}
      </div>
      <div className={styles.tabPanel} role="tabpanel">
        {tabs[safeIndex].props.children}
      </div>
    </div>
  );
}

export { TabsBlock, TabBlock };
export type { TabsBlockProps, TabBlockProps };
