<script lang="ts">
  import type { TouchEventHandler } from "svelte/elements";

  let startX = $state<number | undefined>(undefined);
  let currentX = $state<number | undefined>(undefined);
  let translateX = $state(0);
  let swiping = $state(false);
  let actionTriggered = $state(true);

  const threshold = 100; // Minimum swipe distance to trigger action.

  const handleTouchStart: TouchEventHandler<HTMLElement> = (e) => {
    startX = e.touches[0].clientX;
    swiping = true;
  };

  const handleTouchMove: TouchEventHandler<HTMLElement> = (e) => {
    if (!swiping) return;
    currentX = e.touches[0].clientX;
    let deltaX = currentX - startX!!;
    if (deltaX < 0) {
      // Only allow left swipe
      translateX = deltaX;
    }
  };

  const handleTouchEnd = () => {
    swiping = false;
    if (Math.abs(translateX) > threshold && !actionTriggered) {
      // Trigger action
      actionTriggered = true;
      console.log("Triggered action");
    } else {
      // Reset position
      translateX = 0;
    }
  };

  let actionWidth = $derived(`${translateX * -1}px`);
  $inspect(actionWidth);
</script>

<div class="swipeable-item-container">
  <div class="action" style="--action-width: {actionWidth}">Action</div>
  <div
    class="item-content"
    ontouchstart={handleTouchStart}
    ontouchmove={handleTouchMove}
    ontouchend={handleTouchEnd}
    style="transform: translateX({translateX}px);"
  >
    Hello
  </div>
</div>

<style>
  .swipeable-item-container {
    position: relative;
    overflow: hidden;
  }

  .action {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: var(--action-width);
    background-color: hsl(var(--red-400));
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    cursor: pointer;
    transition: width 0.3s ease;
  }

  .item-content {
    background-color: var(--bg-200);
    position: relative;
    padding: 16px;
    display: flex;
    align-items: center;
    transition: transform 0.3s ease;
    /* Ensure the item content is above the action button */
    z-index: 1;
    touch-action: pan-y;
  }

  /* Optional: Hover effect for desktop */
  @media (hover: hover) {
    .swipeable-item-container:hover .action {
      display: flex;
    }
  }
</style>
