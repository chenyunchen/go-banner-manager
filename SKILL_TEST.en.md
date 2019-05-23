# Task

* A product manager at Mercari has a great idea!  Let’s add banners to the Mercari app.
* We need some way of controlling <u>**when**</u> the banners are displayed.  
* The engineering manager has asked you to encapsulate this behavior in a Banner class.  
* Design the API for the Banner class, and implement the class according to the following specifications.
* Technical Language
    * A banner is **expired** when the display period is over.
    * A banner’s **display period** is the duration the banner is active on the screen. 
    * A banner is **active** during the display period.


# Requirements

* You may code in any of the following languages(We use Go and PHP mainly in the company)
  * Go, PHP, Java
* We’re looking for a well designed API, a clean implementation, and a well designed class structure (if applicable).  
* **This is not a REST API**, assume that this class will be isolated from the controller tier.  If you think of this within a web framework, it’s a class that might be instantiated within the view layer.
* Use the standard library of your chosen programming language as much as possible.  Don’t go inventing your own timezone framework.  
    * If you are actually the inventor of a timezone framework, feel free to submit that as your project instead. :)
* Ensure that the solution is well tested, with appropriate documentation on how to run the tests.
* You also do not have to implement a data layer, feel free to stub this out with mock or test data.  However, the code should be written in a way that is easily adapted to accept a database layer in the future.


# Specifications

* Banner Display Period Conditions
    * Each banner is associated with a promotion.
    * Therefore, each banner will only run for a specific period of time.
    * **Ensure the display period can be set individually for each banner.**
* Banner Display Rules
    * If the banner is **within** the display period, display the banner.
    * The banner display rules should be **timezone aware**.
    * Only one banner can be displayed at a time.
* Internal Release & QA Considerations
* We’d like to display the banner **if the user has an internal IP address** (10.0.0.1, 10.0.0.2), even if the current time is **before** the display period of the banner.
* **After** a banner expires, it should not be displayed again.
* During QA, there may be occasions where two banners are considered active.  In this case, the banner with the earlier expiration should be displayed.

