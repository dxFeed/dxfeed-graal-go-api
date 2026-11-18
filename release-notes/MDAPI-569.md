* **\[MDAPI-569]\[GO]** Fixed `DXFeedSubscription.RemoveListener`: the removed listener no longer receives events, and
  its native listener is released. A listener may remove itself from `Update`. `Close` releases the remaining
  listeners.
