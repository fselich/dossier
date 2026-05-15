## Context

The TUI builds the view in `View()` by assembling rows with `strings.Join`. The current structure is:

```
boxTop
addBorderSides(renderHeader)
addBorderSides(renderTabBar)
[addBorderSides(renderSpecSubnav)]  ← here, before the separator
boxInnerSep
addBorderSides(vp.View)
boxInnerSep
addBorderSides(renderHelpBar)
boxBottom
```

The change consists of moving `renderSpecSubnav` to after the `boxInnerSep`, inside the content block:

```
boxTop
addBorderSides(renderHeader)
addBorderSides(renderTabBar)
boxInnerSep
[addBorderSides(renderSpecSubnav)]  ← here, first line of the content block
addBorderSides(vp.View)
boxInnerSep
addBorderSides(renderHelpBar)
boxBottom
```

## Goals / Non-Goals

**Goals:**
- The spec chips appear visually separated from the tab bar by the `├───┤` separator.
- Keep the viewport height calculation correct (the subnav still subtracts 1 line from the content).

**Non-Goals:**
- Changing styles, colours or navigation logic.
- Altering key behaviour or the cycle between specs.

## Decisions

**Keep the subnav as a static row outside the viewport (do not embed it in the scrollable content)**

Discarded alternative: prepending the chip line to the viewport content. This would cause the chips to disappear on scroll, which would be confusing. The row must be fixed.

The correct solution is to continue rendering the subnav as an additional row in `View()`, simply shifted in position (after `boxInnerSep`). The height adjustment in `contentHeight()` stays the same — the subnav still occupies 1 line of the interior space.

## Risks / Trade-offs

- [Minimal risk] The only perceptible behavioural change is visual (position of the subnav). There are no changes to logic, state or keyboard handling.
